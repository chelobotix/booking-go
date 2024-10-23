package repository

import "github.com/chelobotix/booking-go/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) error
}
