package usecase

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/evrintobing17/JobVacancy/app/helpers/jwthelper"
	"github.com/evrintobing17/JobVacancy/app/models"
	"github.com/evrintobing17/JobVacancy/app/modules/job"
	"github.com/jinzhu/gorm"
)

type UC struct {
	repo job.JobRepository
}

var (
	JWTDuration = 1 * time.Hour

	PasswordHashCost = 14

	//ErrInvalidCredential - Standard error message for invalid credential
	ErrInvalidCredential    = errors.New("invalid credential")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUserNotFound         = errors.New("user not found")
	ErrUIDNotYetRegistered  = errors.New("the user id is not yet registered to DB")
	ErrFirebaseTokenInvalid = errors.New("firebase authentication token is invalid")
	ErrGameKeyInvalid       = errors.New("error game key invalid")
	ErrEmailAlreadyExist    = errors.New("email already exist")
	ErrUsernameAlreadyExist = errors.New("username already exist")
	ErrPhoneAlreadyExist    = errors.New("phone already exist")
	ErrURLCallbackNotSet    = errors.New("callback url not set")
	//time in second
	minute int64 = 60
	hour         = minute * 60
	day          = 24 * hour
)

func NewUserUsecase(repo job.JobRepository) job.JobUsecase {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Login(email, password string) (user *models.User, token string, err error) {
	user, err = uc.repo.GetByUsername(email)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, "", ErrInvalidCredential
		}
		return nil, "", err
	}

	//compare password against Hash

	if user.Password != password {
		return nil, "", ErrInvalidCredential
	}

	jwtExpirationDurationDayString := os.Getenv("jwt.expirationDurationDay")
	var jwtExpirationDurationDay int
	jwtExpirationDurationDay, err = strconv.Atoi(jwtExpirationDurationDayString)
	if err != nil {
		return nil, "", err
	}

	// Conversion to seconds
	jwtExpiredAt := time.Now().Unix() + int64(jwtExpirationDurationDay*3600*24)

	userClaims := jwthelper.AccessJWTClaims{Id: user.ID, ExpiresAt: jwtExpiredAt}
	jwtToken, err := jwthelper.NewWithClaims(userClaims)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}

func (uc *UC) RefreshAccessJWT(userId int) (newAccessJWT string, err error) {
	//create new AccessJWT
	accessJWT, err := uc.generateUserJWTDriver(userId)
	if err != nil {
		return "", err
	}

	return accessJWT, nil
}

func (uc *UC) generateUserJWTDriver(userId int) (token string, err error) {
	//Create JWT
	jwtExpiredAt := time.Now().Add(JWTDuration).Unix()

	userClaims := jwthelper.AccessJWTClaims{
		Id: userId, ExpiresAt: jwtExpiredAt,
	}
	jwtToken, err := jwthelper.NewWithClaims(userClaims)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
