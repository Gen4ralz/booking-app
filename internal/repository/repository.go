package repository

import "github.com/gen4ralz/booking-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) error
}