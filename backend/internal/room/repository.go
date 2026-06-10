package room

import (
	"github.com/abram056/syncstream/backend/internal/models"
)

type Repository interface {
	CreateRoom(room *models.Room) error
	GetRoomByID(id string) (*models.Room, error)
	DeleteRoom(id string) error
}
