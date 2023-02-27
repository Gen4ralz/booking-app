package dbrepo

import (
	"errors"
	"time"

	"github.com/gen4ralz/booking-app/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDates return true if available exists for roomID, and false if no available
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool,error) {
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

func (m *testDBRepo) GetRoomById(id int) (models.Room, error){
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}

	return room,nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error){
	var user models.User
	return user,nil
}

func (m *testDBRepo) UpdateUser(user models.User) error{
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string)(int,string,error){
	return 1,"",nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}