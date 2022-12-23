package services

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/xmas-faf-api/dto"
	"github.com/Marcel-MD/xmas-faf-api/logger"
	"github.com/Marcel-MD/xmas-faf-api/models"
	"github.com/Marcel-MD/xmas-faf-api/repositories"
	"github.com/rs/zerolog/log"
)

type ITrainingService interface {
	FindAll() []models.Training
	FindOne(id string) (models.Training, error)
	Create(dto dto.CreateTraining, userID string) (models.Training, error)
	Update(trainingID, userID string, dto dto.UpdateTraining) (models.Training, error)
	Delete(trainingID, userID string) error
	AddUser(trainingID, addUserID, userID string) error
	RemoveUser(trainingID, removeUserID, userID string) error
	VerifyUserInTraining(trainingID, userID string) error
}

type TrainingService struct {
	trainingRepository repositories.ITrainingRepository
	userRepository     repositories.IUserRepository
}

var (
	trainingOnce    sync.Once
	trainingService ITrainingService
)

func GetTrainingService() ITrainingService {
	trainingOnce.Do(func() {
		log.Info().Msg("Initializing training service")
		trainingService = &TrainingService{
			trainingRepository: repositories.GetTrainingRepository(),
			userRepository:     repositories.GetUserRepository(),
		}
	})
	return trainingService
}

func (s *TrainingService) FindAll() []models.Training {
	log.Debug().Msg("Finding all trainings")

	return s.trainingRepository.FindAll()
}

func (s *TrainingService) FindOne(id string) (models.Training, error) {
	log.Debug().Str(logger.TrainingID, id).Msg("Finding training")

	training, err := s.trainingRepository.FindByIdWithUsers(id)
	if err != nil {
		return training, err
	}

	return training, nil
}

func (s *TrainingService) Create(dto dto.CreateTraining, userID string) (models.Training, error) {
	log.Debug().Str(logger.UserID, userID).Msg("Creating training")

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return models.Training{}, err
	}

	training := models.Training{
		Name:    dto.Name,
		OwnerID: userID,
	}

	err = s.trainingRepository.Create(&training)
	if err != nil {
		return training, err
	}

	err = s.AddUser(training.ID, user.ID, userID)
	if err != nil {
		return training, err
	}

	return training, nil
}

func (s *TrainingService) Update(trainingID, userID string, dto dto.UpdateTraining) (models.Training, error) {
	log.Debug().Str(logger.TrainingID, trainingID).Str(logger.UserID, userID).Msg("Updating training")

	training, err := s.trainingRepository.FindByID(trainingID)
	if err != nil {
		return training, err
	}

	if training.OwnerID != userID {
		return training, errors.New("you are not the owner of this training")
	}

	training.Name = dto.Name

	err = s.trainingRepository.Update(&training)
	if err != nil {
		return training, err
	}

	return training, nil
}

func (s *TrainingService) Delete(trainingID, userID string) error {
	log.Debug().Str(logger.TrainingID, trainingID).Str(logger.UserID, userID).Msg("Deleting training")

	training, err := s.trainingRepository.FindByID(trainingID)
	if err != nil {
		return err
	}

	if training.OwnerID != userID {
		return errors.New("you are not the owner of this training")
	}

	err = s.trainingRepository.Delete(&training)
	if err != nil {
		return err
	}

	return nil
}

func (s *TrainingService) AddUser(trainingID, addUserID, userID string) error {
	log.Debug().Str(logger.TrainingID, trainingID).Msg("Adding user to training")

	err := s.trainingRepository.VerifyUserInTraining(trainingID, addUserID)
	if err == nil {
		return errors.New("user already in this training")
	}

	training, err := s.trainingRepository.FindByID(trainingID)
	if err != nil {
		return err
	}

	if training.OwnerID != userID {
		return errors.New("you are not the owner of this training")
	}

	user, err := s.userRepository.FindByID(addUserID)
	if err != nil {
		return err
	}

	err = s.trainingRepository.AddUser(&training, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *TrainingService) RemoveUser(trainingID, removeUserID, userID string) error {
	log.Debug().Str(logger.TrainingID, trainingID).Str(logger.UserID, removeUserID).Msg("Removing user from training")

	err := s.trainingRepository.VerifyUserInTraining(trainingID, removeUserID)
	if err != nil {
		return err
	}

	training, err := s.trainingRepository.FindByID(trainingID)
	if err != nil {
		return err
	}

	user, err := s.userRepository.FindByID(removeUserID)
	if err != nil {
		return err
	}

	if training.OwnerID != userID && removeUserID != userID {
		return errors.New("unauthorized")
	}

	if training.OwnerID == removeUserID {
		return errors.New("you are the owner of this training")
	}

	err = s.trainingRepository.RemoveUser(&training, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *TrainingService) VerifyUserInTraining(trainingID, userID string) error {
	log.Debug().Str(logger.TrainingID, trainingID).Str(logger.UserID, userID).Msg("Verifying user in training")
	return s.trainingRepository.VerifyUserInTraining(trainingID, userID)
}
