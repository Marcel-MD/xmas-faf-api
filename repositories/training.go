package repositories

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/xmas-faf-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ITrainingRepository interface {
	FindAll() []models.Training
	FindByID(id string) (models.Training, error)
	FindByIdWithUsers(id string) (models.Training, error)
	Create(training *models.Training) error
	Update(training *models.Training) error
	Delete(training *models.Training) error
	AddUser(training *models.Training, user *models.User) error
	RemoveUser(training *models.Training, user *models.User) error
	VerifyUserInTraining(trainingID, userID string) error
}

type TrainingRepository struct {
	DB *gorm.DB
}

var (
	trainingOnce       sync.Once
	trainingRepository ITrainingRepository
)

func GetTrainingRepository() ITrainingRepository {
	trainingOnce.Do(func() {
		log.Info().Msg("Initializing training repository")
		trainingRepository = &TrainingRepository{
			DB: models.GetDB(),
		}
	})
	return trainingRepository
}

func (r *TrainingRepository) FindAll() []models.Training {
	var trainings []models.Training
	r.DB.Find(&trainings).Preload("Users")
	return trainings
}

func (r *TrainingRepository) FindByID(id string) (models.Training, error) {
	var training models.Training
	err := r.DB.First(&training, "id = ?", id).Preload("Users").Preload("Posts.Comments").Preload("Posts.Files").Error

	return training, err
}

func (r *TrainingRepository) FindByIdWithUsers(id string) (models.Training, error) {
	var training models.Training
	err := r.DB.Model(&models.Training{}).Preload("Users").Preload("Posts.Comments").Preload("Posts.Files").First(&training, "id = ?", id).Error

	return training, err
}

func (r *TrainingRepository) Create(training *models.Training) error {
	return r.DB.Create(training).Error
}

func (r *TrainingRepository) Update(training *models.Training) error {
	return r.DB.Save(training).Error
}

func (r *TrainingRepository) Delete(training *models.Training) error {
	return r.DB.Delete(training).Error
}

func (r *TrainingRepository) AddUser(training *models.Training, user *models.User) error {
	return r.DB.Model(training).Omit("Users.*").Association("Users").Append(user)
}

func (r *TrainingRepository) RemoveUser(training *models.Training, user *models.User) error {
	return r.DB.Model(training).Association("Users").Delete(user)
}

func (r *TrainingRepository) VerifyUserInTraining(trainingID, userID string) error {
	var training models.Training
	err := r.DB.Model(&models.Training{}).Preload("Users").First(&training, "id = ?", trainingID).Error
	if err != nil {
		return err
	}

	for _, user := range training.Users {
		if user.ID == userID {
			return nil
		}
	}

	return errors.New("user is not in training")
}
