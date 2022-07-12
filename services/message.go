package services

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type IMessageService interface {
	FindByRoomID(roomID, userID string, params dto.MessageQueryParams) ([]models.Message, error)
	Create(roomID, userID string, dto dto.CreateMessage) (models.Message, error)
	Update(id, userID string, dto dto.UpdateMessage) (models.Message, error)
	Delete(id, userID string) error
	VerifyUserInRoom(roomID, userID string) error
}

type MessageService struct {
	DB *gorm.DB
}

var (
	messageOnce    sync.Once
	messageService IMessageService
)

func GetMessageService() IMessageService {
	messageOnce.Do(func() {
		log.Info().Msg("Initializing message service")
		messageService = &MessageService{
			DB: models.GetDB(),
		}
	})
	return messageService
}

func (s *MessageService) FindByRoomID(roomID, userID string, params dto.MessageQueryParams) ([]models.Message, error) {
	log.Debug().Str("room_id", roomID).Str("user_id", userID).Msg("Finding messages")

	var messages []models.Message

	err := s.VerifyUserInRoom(roomID, userID)
	if err != nil {
		return messages, err
	}

	s.DB.Scopes(models.Paginate(params.Page, params.Size)).Model(&models.Message{}).
		Order("created_at desc").Preload("User").Find(&messages, "room_id = ?", roomID)

	return messages, nil
}

func (s *MessageService) Create(roomID, userID string, dto dto.CreateMessage) (models.Message, error) {
	log.Debug().Str("room_id", roomID).Str("user_id", userID).Msg("Creating message")

	var message models.Message
	err := s.VerifyUserInRoom(roomID, userID)
	if err != nil {
		return message, err
	}

	var user models.User
	err = s.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		return message, err
	}

	message.Text = dto.Text
	message.RoomID = roomID
	message.UserID = userID

	err = s.DB.Create(&message).Error
	if err != nil {
		return message, err
	}

	message.User = user

	return message, nil
}

func (s *MessageService) Update(id, userID string, dto dto.UpdateMessage) (models.Message, error) {
	log.Debug().Str("id", id).Str("user_id", userID).Msg("Updating message")

	var message models.Message
	err := s.DB.First(&message, "id = ?", id).Error
	if err != nil {
		return message, err
	}

	if message.UserID != userID {
		return message, errors.New("you are not allowed to update this message")
	}

	message.Text = dto.Text

	err = s.DB.Save(&message).Error
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) Delete(id, userID string) error {
	log.Debug().Str("id", id).Str("user_id", userID).Msg("Deleting message")

	var message models.Message
	err := s.DB.First(&message, "id = ?", id).Error
	if err != nil {
		return err
	}

	if message.UserID != userID {
		return errors.New("you are not allowed to delete this message")
	}

	err = s.DB.Delete(&message).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) VerifyUserInRoom(roomID, userID string) error {
	log.Debug().Str("room_id", roomID).Str("user_id", userID).Msg("Verifying user in room")

	var room models.Room
	err := s.DB.Model(&models.Room{}).Preload("Users").First(&room, "id = ?", roomID).Error
	if err != nil {
		return err
	}

	for _, user := range room.Users {
		if user.ID == userID {
			return nil
		}
	}

	return errors.New("user is not in room")
}