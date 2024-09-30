package repository

import (
	"github.com/Searedphantasm/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(rr models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)
	Authenticate(email, testPassword string) (int, string, error)
	UpdateUser(u models.User) error

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(res models.Reservation) error
	DeleteReservationByID(id int) error
	UpdateProcessedForReservation(id, processed int) error
}
