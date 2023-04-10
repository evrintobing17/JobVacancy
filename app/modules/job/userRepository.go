package job

import "github.com/evrintobing17/JobVacancy/app/models"

type JobRepository interface {
	GetByUsername(username string) (*models.User, error)
}
