package repositories

import (
	"sync"

	"github.com/Marcel-MD/xmas-faf-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	FindByPostID(postID string) []models.Comment
	FindByID(commentId string) (models.Comment, error)
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(comment *models.Comment) error
}

type CommentRepository struct {
	DB *gorm.DB
}

var (
	commentOnce       sync.Once
	commentRepository ICommentRepository
)

func GetCommentRepository() ICommentRepository {
	commentOnce.Do(func() {
		log.Info().Msg("Initializing comment repository")
		commentRepository = &CommentRepository{
			DB: models.GetDB(),
		}
	})
	return commentRepository
}

func (r *CommentRepository) FindByPostID(postID string) []models.Comment {
	var comments []models.Comment

	r.DB.Model(&models.Comment{}).Order("created_at desc").Find(&comments, "post_id = ?", postID)

	return comments
}

func (r *CommentRepository) FindByID(commentId string) (models.Comment, error) {
	var comment models.Comment
	err := r.DB.First(&comment, "id = ?", commentId).Error

	return comment, err
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.DB.Save(comment).Error
}

func (r *CommentRepository) Delete(comment *models.Comment) error {
	return r.DB.Delete(comment).Error
}
