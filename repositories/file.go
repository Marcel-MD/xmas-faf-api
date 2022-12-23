package repositories

import (
	"sync"

	"github.com/Marcel-MD/xmas-faf-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type IFileRepository interface {
	FindByPostID(postID string) []models.File
	FindByID(id string) (models.File, error)
	Create(file *models.File) error
	Update(file *models.File) error
	Delete(file *models.File) error
}

type FileRepository struct {
	DB *gorm.DB
}

var (
	fileOnce       sync.Once
	fileRepository IFileRepository
)

func GetFileRepository() IFileRepository {
	fileOnce.Do(func() {
		log.Info().Msg("Initializing file repository")
		fileRepository = &FileRepository{
			DB: models.GetDB(),
		}
	})
	return fileRepository
}

func (r *FileRepository) FindByPostID(postID string) []models.File {
	var files []models.File

	r.DB.Find(&files, "post_id = ?", postID)

	return files
}

func (r *FileRepository) FindByID(id string) (models.File, error) {
	var file models.File
	err := r.DB.First(&file, "id = ?", id).Error

	return file, err
}

func (r *FileRepository) Create(file *models.File) error {
	return r.DB.Create(file).Error
}

func (r *FileRepository) Update(file *models.File) error {
	return r.DB.Save(file).Error
}

func (r *FileRepository) Delete(file *models.File) error {
	return r.DB.Delete(file).Error
}
