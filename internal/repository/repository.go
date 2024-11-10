package repository

import (
	"github.com/chelobotix/booking-go/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(r models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDateByRoomId(startDate, endDate time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
	GetRestrictionsForRoomByDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error)
	GetAllRooms() ([]models.Room, error)
	GetUserById(id int) (models.User, error)
	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservation(id int) (models.Reservation, error)
	UpdateReservation(res models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
}
