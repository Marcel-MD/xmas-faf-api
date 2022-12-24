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

type IPostService interface {
	FindByTrainingID(trainingID, userID string) ([]models.Post, error)
	Create(trainingID, userID string, dto dto.CreatePost) (models.Post, error)
	Update(postID, userID string, dto dto.UpdatePost) (models.Post, error)
	Delete(postID, userID string) (models.Post, error)
}

type PostService struct {
	postRepository     repositories.IPostRepository
	trainingRepository repositories.ITrainingRepository
	userRepository     repositories.IUserRepository
}

var (
	postOnce    sync.Once
	postService IPostService
)

func GetPostService() IPostService {
	postOnce.Do(func() {
		log.Info().Msg("Initializing post service")
		postService = &PostService{
			postRepository:     repositories.GetPostRepository(),
			trainingRepository: repositories.GetTrainingRepository(),
			userRepository:     repositories.GetUserRepository(),
		}
	})
	return postService
}

func (s *PostService) FindByTrainingID(trainingID, userID string) ([]models.Post, error) {
	log.Debug().Str(logger.TrainingID, trainingID).Str(logger.UserID, userID).Msg("Finding posts")

	var posts []models.Post

	training, err := s.trainingRepository.FindByID(trainingID)
	if err != nil {
		return posts, err
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return posts, err
	}

	err = s.verifyIfCanRead(training, user)
	if err != nil {
		return posts, err
	}

	posts = s.postRepository.FindByTrainingID(trainingID)

	return posts, nil
}

func (s *PostService) Create(trainingID, userID string, dto dto.CreatePost) (models.Post, error) {
	log.Debug().Str(logger.TrainingID, trainingID).Str(logger.UserID, userID).Msg("Creating post")

	var post models.Post

	training, err := s.trainingRepository.FindByID(trainingID)
	if err != nil {
		return post, err
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return post, err
	}

	err = s.verifyIfCanWrite(training, user)
	if err != nil {
		return post, err
	}

	post.Text = dto.Text
	post.Title = dto.Title
	post.TrainingID = trainingID
	post.UserID = userID

	err = s.postRepository.Create(&post)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (s *PostService) Update(postID, userID string, dto dto.UpdatePost) (models.Post, error) {
	log.Debug().Str(logger.PostID, postID).Str(logger.UserID, userID).Msg("Updating post")

	post, err := s.postRepository.FindByID(postID)
	if err != nil {
		return post, err
	}

	if post.UserID != userID {
		return post, errors.New("you are not allowed to update this post")
	}

	post.Text = dto.Text

	err = s.postRepository.Update(&post)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (s *PostService) Delete(postID, userID string) (models.Post, error) {
	log.Debug().Str(logger.PostID, postID).Str(logger.UserID, userID).Msg("Deleting post")

	post, err := s.postRepository.FindByID(postID)
	if err != nil {
		return post, err
	}

	if post.UserID != userID {
		return post, errors.New("you are not allowed to delete this post")
	}

	post.Text = ""

	err = s.postRepository.Update(&post)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (s *PostService) verifyIfCanWrite(training models.Training, user models.User) error {
	log.Debug().Str(logger.TrainingID, training.ID).Str(logger.UserID, user.ID).Msg("Verifying if user is authorized in training")

	if training.OwnerID == user.ID {
		return nil
	}

	return errors.New("you are not allowed to write in this training")
}

func (s *PostService) verifyIfCanRead(training models.Training, user models.User) error {
	log.Debug().Str(logger.TrainingID, training.ID).Str(logger.UserID, user.ID).Msg("Verifying if user is authorized in training")

	return s.trainingRepository.VerifyUserInTraining(training.ID, user.ID)
}
