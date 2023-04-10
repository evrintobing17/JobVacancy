package job

import (
	"github.com/evrintobing17/JobVacancy/app/models"
)

type JobUsecase interface {
	Login(email, password string) (user *models.User, token string, err error)
}
