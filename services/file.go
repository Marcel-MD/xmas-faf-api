package services

import (
	"path"
	"sync"

	"github.com/Marcel-MD/xmas-faf-api/models"
	"github.com/Marcel-MD/xmas-faf-api/repositories"
)

type IFileService interface {
	FindByPostID(postID string) []models.File
	FindByID(id string) (models.File, error)
	Create(postID, fileName string, data []byte) (models.File, error)
	Delete(id string) error
}

type FileService struct {
	blobService    IBlobService
	fileRepository repositories.IFileRepository
}

var (
	fileOnce    sync.Once
	fileService IFileService
)

func GetFileService() IFileService {
	fileOnce.Do(func() {
		fileService = &FileService{
			blobService:    GetBlobService(),
			fileRepository: repositories.GetFileRepository(),
		}
	})

	return fileService
}

func (s *FileService) FindByPostID(postID string) []models.File {
	return s.fileRepository.FindByPostID(postID)
}

func (s *FileService) FindByID(id string) (models.File, error) {
	return s.fileRepository.FindByID(id)
}

func (s *FileService) Create(postID, fileName string, data []byte) (models.File, error) {

	url, err := s.blobService.Upload(fileName, data)
	if err != nil {
		return models.File{}, err
	}

	file := models.File{
		PostID: postID,
		Name:   fileName,
		Ext:    path.Ext(fileName),
		Url:    url,
	}

	err = s.fileRepository.Create(&file)
	if err != nil {
		return file, err
	}

	return file, nil
}

func (s *FileService) Delete(id string) error {
	file, err := s.fileRepository.FindByID(id)
	if err != nil {
		return err
	}

	err = s.blobService.Delete(file.Name)
	if err != nil {
		return err
	}

	return s.fileRepository.Delete(&file)
}
