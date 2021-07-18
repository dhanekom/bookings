package dbrepo

import (
	"context"
	"time"

	"github.com/dhanekom/bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var newId int

	stmt := `insert into reservations (first_name, last_name, email, phone, 
		       start_date, end_date, room_id, created_at, updated_at)
	         values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
		       created_at, updated_at, restriction_id)
					 values
					 ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID else returns false
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `
	select count(id)
	from room_restrictions
	where room_id = $1
	  and $2 < end_date and $3 > start_date`

	var numRows int
	err := m.DB.QueryRowContext(ctx, stmt, roomID, start, end).Scan(&numRows)
	if err != nil {
		return false, err
	}
	return numRows == 0, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for a given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `
	select r.id, r.room_name
	from rooms r
	where not r.id in (select room_id
										 from room_restrictions rr
										 where $1 < end_date and $2 > start_date)`

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, nil
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.RoomName,
		&room.CreateAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}
