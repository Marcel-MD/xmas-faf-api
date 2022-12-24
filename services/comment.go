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

type ICommentService interface {
	FindByPostID(postID, userID string) ([]models.Comment, error)
	Create(postID, userID string, dto dto.CreateComment) (models.Comment, error)
	Update(commentID, userID string, dto dto.UpdateComment) (models.Comment, error)
	Delete(commentID, userID string) (models.Comment, error)
}

type CommentService struct {
	commentRepository  repositories.ICommentRepository
	postRepository     repositories.IPostRepository
	trainingRepository repositories.ITrainingRepository
	userRepository     repositories.IUserRepository
}

var (
	commentOnce    sync.Once
	commentService ICommentService
)

func GetCommentService() ICommentService {
	commentOnce.Do(func() {
		log.Info().Msg("Initializing comment service")
		commentService = &CommentService{
			commentRepository:  repositories.GetCommentRepository(),
			postRepository:     repositories.GetPostRepository(),
			trainingRepository: repositories.GetTrainingRepository(),
			userRepository:     repositories.GetUserRepository(),
		}
	})
	return commentService
}

func (s *CommentService) FindByPostID(postID, userID string) ([]models.Comment, error) {
	log.Debug().Str(logger.PostID, postID).Str(logger.UserID, userID).Msg("Finding comments")

	var comments []models.Comment

	_, err := s.postRepository.FindByID(postID)
	if err != nil {
		return comments, err
	}

	_, err = s.userRepository.FindByID(userID)
	if err != nil {
		return comments, err
	}

	// err = s.verifyIfCanRead(post.Training, user)
	// if err != nil {
	// 	return comments, err
	// }

	comments = s.commentRepository.FindByPostID(postID)

	return comments, nil
}

func (s *CommentService) Create(postID, userID string, dto dto.CreateComment) (models.Comment, error) {
	log.Debug().Str(logger.PostID, postID).Str(logger.UserID, userID).Msg("Creating comment")

	var comment models.Comment

	_, err := s.postRepository.FindByID(postID)
	if err != nil {
		return comment, err
	}

	_, err = s.userRepository.FindByID(userID)
	if err != nil {
		return comment, err
	}

	// err = s.verifyIfCanWrite(post.Training, user)
	// if err != nil {
	// 	return comment, err
	// }

	comment.Text = dto.Text
	comment.PostID = postID
	comment.UserID = userID

	err = s.commentRepository.Create(&comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (s *CommentService) Update(commentID, userID string, dto dto.UpdateComment) (models.Comment, error) {
	log.Debug().Str(logger.CommentId, commentID).Str(logger.UserID, userID).Msg("Updating comment")

	comment, err := s.commentRepository.FindByID(commentID)
	if err != nil {
		return comment, err
	}

	if comment.UserID != userID {
		return comment, errors.New("you are not allowed to update this comment")
	}

	comment.Text = dto.Text

	err = s.commentRepository.Update(&comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (s *CommentService) Delete(commentID, userID string) (models.Comment, error) {
	log.Debug().Str(logger.CommentId, commentID).Str(logger.UserID, userID).Msg("Deleting comment")

	comment, err := s.commentRepository.FindByID(commentID)
	if err != nil {
		return comment, err
	}

	if comment.UserID != userID {
		return comment, errors.New("you are not allowed to delete this comment")
	}

	comment.Text = ""

	err = s.commentRepository.Update(&comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// func (s *CommentService) verifyIfCanWrite(training models.Training, user models.User) error {
// 	log.Debug().Str(logger.TrainingID, training.ID).Str(logger.UserID, user.ID).Msg("Verifying if user is authorized in training")

// 	if training.OwnerID == user.ID {
// 		return nil
// 	}

// 	return errors.New("you are not allowed to write in this training")
// }

// func (s *CommentService) verifyIfCanRead(training models.Training, user models.User) error {
// 	log.Debug().Str(logger.TrainingID, training.ID).Str(logger.UserID, user.ID).Msg("Verifying if user is authorized in training")

// 	return s.trainingRepository.VerifyUserInTraining(training.ID, user.ID)
// }
