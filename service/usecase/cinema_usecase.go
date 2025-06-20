package usecase

import (
	"cinema-seat-reservation/service/cinema_service"
	"cinema-seat-reservation/service/model/request"
	"cinema-seat-reservation/service/model/response"
	"errors"
	"sync"
	"sync/atomic"
)

type CinemaServiceUseCase interface {
	Configure(req request.ConfigRequest)
	QueryAvailableSeats(req request.QueryAvailableSeatRequest) response.QueryAvailableSeatResp
	CheckAvailableSeats(req request.CheckAvailableSeatRequest) response.CheckAvailableSeatResp
	ReserveSeats(req request.ReserveSeatRequest) error
	CancelSeats(req request.CancelSeatRequest) error
}

type cinemaServiceUsecase struct {
	cs *cinema_service.CinemaService
}

func NewCinemaUsecase(cinemaService *cinema_service.CinemaService) CinemaServiceUseCase {
	return &cinemaServiceUsecase{
		cs: cinemaService,
	}
}

// Configure Cinema Layout
func (s *cinemaServiceUsecase) Configure(req request.ConfigRequest) {
	s.cs.Configure(req.Rows, req.Columns, req.MinDistance)
}

// Query Available Seats
func (s *cinemaServiceUsecase) QueryAvailableSeats(req request.QueryAvailableSeatRequest) (result response.QueryAvailableSeatResp) {

	count := req.Count
	for i := 0; i < s.cs.Rows; i++ {
		for j := 0; j <= s.cs.Cols-count; j++ {
			valid := true
			block := make([]response.Seat, count)
			for k := 0; k < count; k++ {
				seat := s.cs.Seats[i][j+k]
				if seat.Taken || !s.validDistance(i, j+k) {
					valid = false
					break
				}
				seatReq := *seat

				block[k] = response.Seat{
					Row:    seatReq.Row,
					Column: seatReq.Column,
					Taken:  seatReq.Taken,
					Group:  seatReq.Group,
				}
			}
			if valid {
				result.AvailableSeats = append(result.AvailableSeats, block)
			}
		}
	}
	return result
}

// Check Available Seats
func (s *cinemaServiceUsecase) CheckAvailableSeats(req request.CheckAvailableSeatRequest) (available response.CheckAvailableSeatResp) {

	seats := req.Seats

	for _, c := range seats {
		if c.Row < s.cs.Rows && c.Column < s.cs.Cols {
			//&& !s.cs.Seats[c.Row][c.Column].Taken {
			lock := s.cs.SeatLocks[c.Row][c.Column]
			lock.Lock()
			if !s.cs.Seats[c.Row][c.Column].Taken {
				seat := *s.cs.Seats[c.Row][c.Column]
				available.AvailableSeats = append(available.AvailableSeats, response.Seat{
					Row:    seat.Row,
					Column: seat.Column,
					Taken:  seat.Taken,
					Group:  seat.Group,
				})
			}
			lock.Unlock()
		}
	}
	return available
}

// Reserve Seats
func (s *cinemaServiceUsecase) ReserveSeats(req request.ReserveSeatRequest) error {
	seats := req.Seats
	locks := make([]*sync.Mutex, len(seats))
	// Lock all seats in order
	for i, c := range seats {
		if c.Row >= s.cs.Rows || c.Column >= s.cs.Cols {
			return errors.New("seat out of bounds")
		}
		locks[i] = s.cs.SeatLocks[c.Row][c.Column]
	}

	for i := range locks {
		locks[i].Lock()
	}

	defer func() {
		for i := range locks {
			locks[i].Unlock()
		}
	}()

	for _, c := range seats {
		seat := s.cs.Seats[c.Row][c.Column]
		if seat.Taken || !s.validDistance(c.Row, c.Column) {
			return errors.New("seat unavailable or violates distance")
		}
	}

	groupID := atomic.AddUint64(&s.cs.GroupCounter, 1)

	for _, c := range seats {
		seat := s.cs.Seats[c.Row][c.Column]
		seat.Taken = true
		seat.Group = groupID
	}

	return nil
}

// Cancel Reservation
func (s *cinemaServiceUsecase) CancelSeats(req request.CancelSeatRequest) error {

	coords := req.Seats
	for _, c := range coords {
		if c.Row < s.cs.Rows && c.Column < s.cs.Cols {
			s.cs.SeatLocks[c.Row][c.Column].Lock()
			s.cs.Seats[c.Row][c.Column].Taken = false
			s.cs.Seats[c.Row][c.Column].Group = 0
			s.cs.SeatLocks[c.Row][c.Column].Unlock()
		}
	}

	return nil
}

// helpers

func (s *cinemaServiceUsecase) validDistance(row, col int) bool {
	for i := 0; i < s.cs.Rows; i++ {
		for j := 0; j < s.cs.Cols; j++ {
			seat := s.cs.Seats[i][j]
			if seat.Taken && s.manhattan(seat.Row, seat.Column, row, col) < s.cs.MinDistance {
				return false
			}
		}
	}
	return true
}

func (s *cinemaServiceUsecase) manhattan(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
