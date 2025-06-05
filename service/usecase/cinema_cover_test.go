package usecase

import (
	"github.com/stretchr/testify/assert"
	"home-assignment/home-assignments/service/cinema_service"
	"home-assignment/home-assignments/service/model/request"
	"testing"
)

func TestCinemaService_AllSeatsReserved(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)

	svc.Configure(request.ConfigRequest{
		Rows:        1,
		Columns:     2,
		MinDistance: 1,
	})

	req := request.ReserveSeatRequest{
		Seats: []request.Seat{{Row: 0, Column: 0}, {Row: 0, Column: 1}},
	}
	err := svc.ReserveSeats(req)
	assert.NoError(t, err)

	// Try to reserve again
	err = svc.ReserveSeats(req)
	assert.Error(t, err, "Should fail because all seats are reserved")
}

func TestCinemaService_OverlappingRequests(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)

	svc.Configure(request.ConfigRequest{
		Rows:        1,
		Columns:     3,
		MinDistance: 1,
	})

	first := request.ReserveSeatRequest{
		Seats: []request.Seat{{Row: 0, Column: 0}, {Row: 0, Column: 1}},
	}

	second := request.ReserveSeatRequest{
		Seats: []request.Seat{{Row: 0, Column: 1}, {Row: 0, Column: 2}},
	}

	err1 := svc.ReserveSeats(first)
	err2 := svc.ReserveSeats(second)

	assert.NoError(t, err1)
	assert.Error(t, err2, "Overlapping seat should fail")
}

func TestCinemaService_InvalidCoordinates(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)

	svc.Configure(request.ConfigRequest{
		Rows:        2,
		Columns:     2,
		MinDistance: 1,
	})

	invalid := request.ReserveSeatRequest{
		Seats: []request.Seat{{Row: 5, Column: 5}},
	}

	err := svc.ReserveSeats(invalid)
	assert.Error(t, err, "Should fail due to out-of-bounds seat")
}

func TestCinemaService_TooManySeatsRequested(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)

	svc.Configure(request.ConfigRequest{
		Rows:        2,
		Columns:     1,
		MinDistance: 2,
	})

	tooMany := request.ReserveSeatRequest{
		Seats: []request.Seat{{Row: 0, Column: 0}, {Row: 0, Column: 1}, {Row: 0, Column: 2}},
	}

	err := svc.ReserveSeats(tooMany)
	assert.Error(t, err, "Should fail because too many seats are requested")
}
