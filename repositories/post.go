package repositories

import (
	"sync"

	"github.com/Marcel-MD/xmas-faf-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type IPostRepository interface {
	FindByTrainingID(trainingID string) []models.Post
	FindByID(id string) (models.Post, error)
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(post *models.Post) error
}

type PostRepository struct {
	DB *gorm.DB
}

var (
	postOnce       sync.Once
	postRepository IPostRepository
)

func GetPostRepository() IPostRepository {
	postOnce.Do(func() {
		log.Info().Msg("Initializing post repository")
		postRepository = &PostRepository{
			DB: models.GetDB(),
		}
	})
	return postRepository
}

func (r *PostRepository) FindByTrainingID(trainingID string) []models.Post {
	var posts []models.Post

	r.DB.Model(&models.Post{}).Preload("Files").Order("created_at desc").Find(&posts, "training_id = ?", trainingID)

	return posts
}

func (r *PostRepository) FindByID(id string) (models.Post, error) {
	var post models.Post
	err := r.DB.First(&post, "id = ?", id).Preload("Files").Error

	return post, err
}

func (r *PostRepository) Create(post *models.Post) error {
	return r.DB.Create(post).Error
}

func (r *PostRepository) Update(post *models.Post) error {
	return r.DB.Save(post).Error
}

func (r *PostRepository) Delete(post *models.Post) error {
	return r.DB.Delete(post).Error
}
