package usecase

import (
	"cinema-seat-reservation/service/cinema_service"
	"cinema-seat-reservation/service/model/request"
	"cinema-seat-reservation/service/model/response"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

type CinemaServiceUseCase interface {
	Configure(req request.ConfigRequest)
	QueryAvailableSeats(req request.QueryAvailableSeatRequest) response.QueryAvailableSeatResp
	CheckAvailableSeats(req request.CheckAvailableSeatRequest) response.CheckAvailableSeatResp
	ReserveSeats(req request.ReserveSeatRequest) (response.ReserveSeatsResp, error)
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
func (c *cinemaServiceUsecase) CheckAvailableSeats(req request.CheckAvailableSeatRequest) (result response.CheckAvailableSeatResp) {

	seatsCheck := req.Seats
	seatLocks := make([]*sync.Mutex, 0, len(seatsCheck))
	for _, seat := range seatsCheck {
		if seat.Column > c.cs.Cols || seat.Row > c.cs.Cols {
			continue
		}
		seatLock := c.cs.SeatLocks[seat.Row][seat.Column]
		seatLock.Lock()
		seatLocks = append(seatLocks, seatLock)
	}

	defer func() {
		for _, seatLock := range seatLocks {
			seatLock.Unlock()
		}
	}()

	for _, seat := range seatsCheck {
		if seat.Row < c.cs.Rows && seat.Column < c.cs.Cols {
			seatValid := c.cs.Seats[seat.Row][seat.Column]
			if !seatValid.Taken {
				result.AvailableSeats = append(result.AvailableSeats, response.Seat{
					Row:    seatValid.Row,
					Column: seatValid.Column,
					Taken:  seatValid.Taken,
					Group:  seatValid.Group,
				})
			}
		}
	}

	return result
}

// Reserve Seats
func (s *cinemaServiceUsecase) ReserveSeats(req request.ReserveSeatRequest) (response.ReserveSeatsResp, error) {
	seats := req.Seats
	// Lock all seats in order
	for _, c := range seats {
		if c.Row >= s.cs.Rows || c.Column >= s.cs.Cols {
			return response.ReserveSeatsResp{}, fmt.Errorf("seat at row %d, column %d is out of bounds", c.Row, c.Column)
		}
	}
	// Sort seats by row and column to avoid deadlocks
	sort.Slice(seats, func(i, j int) bool {
		if seats[i].Row == seats[j].Row {
			return seats[i].Column < seats[j].Column
		}
		return seats[i].Row < seats[j].Row
	})

	locks := make([]*sync.Mutex, len(seats))
	for i, c := range seats {
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
			return response.ReserveSeatsResp{}, fmt.Errorf("seat at row %d, column %d is already taken or too close to another seat", c.Row, c.Column)
		}
	}

	groupID := atomic.AddUint64(&s.cs.GroupCounter, 1)

	for _, c := range seats {
		seat := s.cs.Seats[c.Row][c.Column]
		seat.Taken = true
		seat.Group = groupID
	}

	return response.ReserveSeatsResp{
		Message: "Seats reserved successfully",
	}, nil
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
