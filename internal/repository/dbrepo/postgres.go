package dbrepo

import (
	"context"
	"github.com/chelobotix/booking-go/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(r models.Reservation) (int, error) {
	var newId int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			 values ($1, $2, $3 , $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
			 values ($1, $2, $3 , $4, $5, $6, $7 )`

	_, err := m.DB.ExecContext(ctx, query,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchAvailabilityByDateByRoomId(startDate, endDate time.Time, roomId int) (bool, error) {
	var result int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT count(id)
			  FROM room_restrictions
			  WHERE $1 < end_date and $2 > start_date and room_id = $3`

	row := m.DB.QueryRowContext(ctx, query, startDate, endDate, roomId)
	err := row.Scan(&result)
	if err != nil {
		return false, err
	}

	if result == 0 {
		return true, nil
	}

	return false, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	var availableRooms []models.Room
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT r.id, r.room_name
			  FROM rooms r
			  WHERE r.id not in(SELECT room_id
			                  FROM room_restrictions rr
			            	  WHERE $1 < rr.end_date and $2 > start_date)`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return nil, err
		}

		availableRooms = append(availableRooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return availableRooms, nil
}
