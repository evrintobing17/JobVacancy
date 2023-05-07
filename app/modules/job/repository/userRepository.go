package repository

import (
	"github.com/evrintobing17/JobVacancy/app/models"
	"github.com/evrintobing17/JobVacancy/app/modules/job"
	"github.com/jinzhu/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) job.JobRepository {
	return &repo{
		db: db,
	}

}

//Get user data by username
func (r *repo) GetByUsername(username string) (*models.User, error) {
	var user models.User

	db := r.db.First(&user, "username = ?", username)
	if db.Error != nil {
		return nil, db.Error
	}
	return &user, nil
}
