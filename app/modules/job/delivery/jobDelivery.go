package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/evrintobing17/JobVacancy/app/helpers/jsonhttpresponse"
	"github.com/evrintobing17/JobVacancy/app/middlewares/authmiddleware"
	"github.com/evrintobing17/JobVacancy/app/models"
	"github.com/evrintobing17/JobVacancy/app/modules/job"
	jobDTO "github.com/evrintobing17/JobVacancy/app/modules/job/delivery/jobDto"
	userUsecase "github.com/evrintobing17/JobVacancy/app/modules/job/usecase"
	"github.com/gin-gonic/gin"
)

type jobHandler struct {
	jobUC          job.JobUsecase
	authMiddleware authmiddleware.AuthMiddleware
}

func NewAuthHTTPHandler(r *gin.Engine, jobUC job.JobUsecase, authMiddleware authmiddleware.AuthMiddleware) {
	handlers := jobHandler{
		jobUC:          jobUC,
		authMiddleware: authMiddleware,
	}

	authorized := r.Group("/v1/job")
	{
		authorized.POST("/login", handlers.login)
	}
	job := r.Group("/v1/job", handlers.authMiddleware.AuthorizeJWTWithUserContext())
	{
		job.GET("/list", handlers.getJobList)
		job.GET("/:id", handlers.jobdetail)
	}
}

func (handler *jobHandler) login(c *gin.Context) {

	var loginReq jobDTO.ReqLogin

	errBind := c.ShouldBind(&loginReq)
	if errBind != nil {
		jsonhttpresponse.BadRequest(c, "")
		return
	}

	_, jwt, err := handler.jobUC.Login(loginReq.Username, loginReq.Password)
	if err != nil {

		if err == userUsecase.ErrInvalidCredential {
			jsonhttpresponse.Unauthorized(c, jsonhttpresponse.NewFailedResponse(err.Error()))
			return
		}

		jsonhttpresponse.InternalServerError(c, jsonhttpresponse.NewFailedResponse(err.Error()))
		return
	}

	response := jobDTO.ResLogin{
		Jwt: jwt,
	}
	jsonhttpresponse.OK(c, response)
}

func (handler *jobHandler) getJobList(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		jsonhttpresponse.BadRequest(c, jsonhttpresponse.NewFailedResponse(errors.New("invalid token")))
		return
	}

	// Extract query parameters
	description := c.Query("description")
	location := c.Query("location")
	fullTime := c.Query("full_time")

	// Build URL with query parameters
	url := "http://dev3.dansmultipro.co.id/api/recruitment/positions.json?"
	if description != "" {
		url += "description=" + description + "&"
	}
	if location != "" {
		url += "location=" + location + "&"
	}
	if fullTime != "" {
		url += "full_time=" + fullTime + "&"
	}

	// Make HTTP request to the API
	resp, err := http.Get(url)
	if err != nil {
		jsonhttpresponse.InternalServerError(c, err)
		return
	}
	defer resp.Body.Close()

	// Parse response into list of jobs
	var jobs []models.Job
	err = json.NewDecoder(resp.Body).Decode(&jobs)
	if err != nil {
		jsonhttpresponse.InternalServerError(c, err.Error())
		return
	}

	// Paginate results
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	start := (parsePageNumber(page) - 1) * 10
	end := start + 10
	if end > len(jobs) {
		end = len(jobs)
	}
	jobs = jobs[start:end]

	// Send response
	jsonhttpresponse.OK(c, jobs)
	return
}

func (handler *jobHandler) jobdetail(c *gin.Context) {
	id := c.Param("id")
	url := "http://dev3.dansmultipro.co.id/api/recruitment/positions/" + id
	resp, err := http.Get(url)
	if err != nil {
		jsonhttpresponse.InternalServerError(c, err.Error())
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var jobdetail interface{}
		err := json.NewDecoder(resp.Body).Decode(&jobdetail)
		if err != nil {
			jsonhttpresponse.InternalServerError(c, err.Error())
			return
		}

		jsonhttpresponse.OK(c, jobdetail)
		return
	}
	jsonhttpresponse.NotFound(c, err.Error())
}

func parsePageNumber(page string) int {
	if page == "" {
		return 1
	}
	n, err := strconv.Atoi(page)
	if err != nil || n < 1 {
		return 1
	}
	return n
}
