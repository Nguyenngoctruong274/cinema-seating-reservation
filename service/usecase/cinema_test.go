package usecase

import (
	"cinema-seat-reservation/service/cinema_service"
	"cinema-seat-reservation/service/model/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCinemaService_Configure(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)
	req := request.ConfigRequest{
		Rows:        4,
		Columns:     5,
		MinDistance: 2,
	}
	svc.Configure(req)

	assert.Equal(t, 4, cinemaService.Rows)
	assert.Equal(t, 5, cinemaService.Cols)
	assert.Equal(t, 2, cinemaService.MinDistance)
	assert.NotNil(t, cinemaService.Seats)
}

func TestCinemaService_QueryAvailableSeats(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)
	req := request.ConfigRequest{
		Rows:        4,
		Columns:     5,
		MinDistance: 4,
	}

	svc.Configure(req)

	blockResp := svc.QueryAvailableSeats(request.QueryAvailableSeatRequest{Count: 3})

	// Expected at least one available block
	assert.NotEmpty(t, blockResp.AvailableSeats)
	for _, block := range blockResp.AvailableSeats {
		assert.Len(t, block, 3)
		for _, seat := range block {
			assert.False(t, seat.Taken)
		}
	}
}

func TestCinemaService_ReserveSeats(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)
	req := request.ConfigRequest{
		Rows:        4,
		Columns:     5,
		MinDistance: 4,
	}

	svc.Configure(req)

	result := request.ReserveSeatRequest{
		Seats: []request.Seat{
			{Row: 0, Column: 0},
			{Row: 0, Column: 1},
		},
	}
	err := svc.ReserveSeats(result)
	assert.NoError(t, err)

	for _, c := range result.Seats {
		assert.True(t, cinemaService.Seats[c.Row][c.Column].Taken)
	}

}

func TestCinemaService_CheckAvailableSeats(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)
	req := request.ConfigRequest{
		Rows:        2,
		Columns:     2,
		MinDistance: 1,
	}

	svc.Configure(req)

	svc.ReserveSeats(request.ReserveSeatRequest{
		Seats: []request.Seat{
			{Row: 0, Column: 0},
		},
	})

	availableResp := svc.CheckAvailableSeats(request.CheckAvailableSeatRequest{
		Seats: []request.Seat{
			{Row: 0, Column: 0}, {Row: 0, Column: 1},
		},
	})

	assert.Len(t, availableResp.AvailableSeats, 1)
	assert.Equal(t, 0, availableResp.AvailableSeats[0].Row)
	assert.Equal(t, 1, availableResp.AvailableSeats[0].Column)
}

func TestCinemaService_CancelSeats(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)
	req := request.ConfigRequest{
		Rows:        2,
		Columns:     2,
		MinDistance: 1,
	}

	svc.Configure(req)
	coords := []request.Seat{
		{Row: 0, Column: 0},
	}
	svc.ReserveSeats(request.ReserveSeatRequest{
		Seats: coords,
	})

	svc.CancelSeats(request.CancelSeatRequest{
		Seats: coords,
	})

	assert.False(t, cinemaService.Seats[0][0].Taken)
	assert.Equal(t, 0, cinemaService.Seats[0][0].Group)
}
