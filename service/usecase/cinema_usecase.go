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
	CheckAvailableSeats(req request.CheckAvailableSeatRequest) (response.CheckAvailableSeatResp, error)
	ReserveSeats(req request.ReserveSeatRequest) (response.ReserveSeatsResp, error)
	CancelSeats(req request.CancelSeatRequest) error
}

type cinemaServiceUsecase struct {
	cinema *cinema_service.CinemaService
}

func NewCinemaUsecase(cinemaService *cinema_service.CinemaService) CinemaServiceUseCase {
	return &cinemaServiceUsecase{
		cinema: cinemaService,
	}
}

// Configure Cinema Layout
func (c *cinemaServiceUsecase) Configure(req request.ConfigRequest) {
	c.cinema.Configure(req.Rows, req.Columns, req.MinDistance)
}

// Query Available Seats
func (c *cinemaServiceUsecase) QueryAvailableSeats(req request.QueryAvailableSeatRequest) (result response.QueryAvailableSeatResp) {
	seatCheck := req.Count
	bookingSet := make(map[[2]int]bool)
	for i := 0; i < c.cinema.Rows; i++ {
		for j := 0; j <= c.cinema.Cols-seatCheck; j++ {
			valid := true
			seatValid := make([]response.Seat, 0, seatCheck)
			for k := 0; k < seatCheck; k++ {
				row := i
				col := j + k
				seat := c.cinema.Seats[row][col]
				if seat.Taken || !c.validDistance(seat.Row, seat.Column, bookingSet) {
					valid = false
					break
				}
				seatValid = append(seatValid, response.Seat{
					Row:    seat.Row,
					Column: seat.Column,
					Taken:  seat.Taken,
					Group:  seat.Group,
				})
			}
			if valid {
				result.AvailableSeats = append(result.AvailableSeats, seatValid)
			}
		}
	}
	return
}

// Check Available Seats
func (c *cinemaServiceUsecase) CheckAvailableSeats(req request.CheckAvailableSeatRequest) (result response.CheckAvailableSeatResp, err error) {
	seatChecks := req.Seats
	seatLocks := make([]*sync.RWMutex, 0, len(seatChecks))
	mapSeatUnique := make(map[[2]int]bool)
	for _, seat := range seatChecks {
		if seat.Row >= c.cinema.Rows || seat.Column >= c.cinema.Cols {
			return response.CheckAvailableSeatResp{}, fmt.Errorf("seat (%d,%d) out of bounds", seat.Row, seat.Column)
		}
		key := [2]int{seat.Row, seat.Row}
		if mapSeatUnique[key] {
			return response.CheckAvailableSeatResp{}, fmt.Errorf("duplicate seat (%d,%d) in request", seat.Row, seat.Column)
		}
		mapSeatUnique[key] = true
	}

	sort.Slice(seatChecks, func(i, j int) bool {
		if seatChecks[i].Row == seatChecks[j].Row {
			return seatChecks[i].Column < seatChecks[j].Column
		}
		return seatChecks[i].Row < seatChecks[j].Row
	})

	for _, seat := range seatChecks {
		if seat.Row >= c.cinema.Rows || seat.Column >= c.cinema.Cols {
			continue
		}
		seatLock := c.cinema.SeatLocks[seat.Row][seat.Column]
		seatLock.RLock()
		seatLocks = append(seatLocks, seatLock)
		seatValid := c.cinema.Seats[seat.Row][seat.Column]
		if !seatValid.Taken {
			result.AvailableSeats = append(result.AvailableSeats, response.Seat{
				Row:    seatValid.Row,
				Column: seatValid.Column,
				Taken:  false,
				Group:  seatValid.Group,
			})
		}
	}

	//unlock
	defer func() {
		for i := len(seatLocks) - 1; i >= 0; i-- {
			seatLocks[i].RUnlock()
		}
	}()

	return
}

// Reserve Seats
func (c *cinemaServiceUsecase) ReserveSeats(req request.ReserveSeatRequest) (response.ReserveSeatsResp, error) {
	bookSeats := req.Seats
	//1. check validate seat
	mapSeatUnique := make(map[[2]int]bool)

	for _, seat := range bookSeats {
		if seat.Row >= c.cinema.Rows || seat.Column >= c.cinema.Cols {
			return response.ReserveSeatsResp{}, fmt.Errorf("seat (%d,%d) out of bounds", seat.Row, seat.Column)
		}
		key := [2]int{seat.Row, seat.Column}
		if mapSeatUnique[key] {
			return response.ReserveSeatsResp{}, fmt.Errorf("duplicate seat (%d,%d) in requests", seat.Row, seat.Column)
		}
		mapSeatUnique[key] = true
	}

	//2. sort avoid deadlocks
	sort.Slice(bookSeats, func(i, j int) bool {
		if bookSeats[i].Row == bookSeats[j].Row {
			return bookSeats[i].Column < bookSeats[j].Column
		}
		return bookSeats[i].Row < bookSeats[j].Row
	})

	//3. lock
	seatLocks := make([]*sync.RWMutex, 0, len(bookSeats))
	for _, seat := range bookSeats {
		if seat.Row < c.cinema.Rows && seat.Column < c.cinema.Cols {
			seatLock := c.cinema.SeatLocks[seat.Row][seat.Column]
			seatLock.Lock()
			seatLocks = append(seatLocks, seatLock)
		}
	}

	defer func() {
		for i := len(seatLocks) - 1; i >= 0; i-- {
			seatLocks[i].Unlock()
		}
	}()

	//4. Check
	for _, seat := range bookSeats {
		if c.cinema.Seats[seat.Row][seat.Column].Taken || !c.validDistance(seat.Row, seat.Column, mapSeatUnique) {
			return response.ReserveSeatsResp{}, fmt.Errorf("seat (%d,%d) is already taken or to close to another seat", seat.Row, seat.Column)
		}
	}
	//5. order
	groupID := atomic.AddUint64(&c.cinema.GroupCounter, 1)
	for _, seat := range bookSeats {
		seatValid := c.cinema.Seats[seat.Row][seat.Column]
		seatValid.Taken = true
		seatValid.Group = groupID
	}

	return response.ReserveSeatsResp{
		Message: "Seats reserved successfully",
	}, nil
}

// Cancel Reservation
func (c *cinemaServiceUsecase) CancelSeats(req request.CancelSeatRequest) error {
	seatsCancel := req.Seats
	mapSeatUnique := make(map[[2]int]bool)
	//	1.Validate
	for _, seat := range seatsCancel {
		if seat.Row >= c.cinema.Rows || seat.Column >= c.cinema.Cols {
			return fmt.Errorf("seat (%d,%d) out of bounds", seat.Row, seat.Column)
		}
		if mapSeatUnique[[2]int{seat.Row, seat.Column}] {
			return fmt.Errorf("duplicate seat (%d,%d) in requests", seat.Row, seat.Column)
		}
	}
	//2.sort
	sort.Slice(seatsCancel, func(i, j int) bool {
		if seatsCancel[i].Row == seatsCancel[j].Row {
			return seatsCancel[i].Column < seatsCancel[i].Column
		}
		return seatsCancel[i].Row < seatsCancel[j].Row
	})
	//3. lock
	seatLocks := make([]*sync.RWMutex, 0, len(seatsCancel))
	for _, seat := range seatsCancel {
		if seat.Row < c.cinema.Rows && seat.Column < c.cinema.Cols {
			seatLock := c.cinema.SeatLocks[seat.Row][seat.Column]
			seatLock.Lock()
			seatLocks = append(seatLocks, seatLock)
		}
	}

	defer func() {
		for i := len(seatLocks) - 1; i >= 0; i-- {
			seatLocks[i].RUnlock()
		}
	}()

	//4. cancel
	for _, seat := range seatsCancel {
		if seat.Row < c.cinema.Rows && seat.Column < c.cinema.Cols {
			seatValid := c.cinema.Seats[seat.Row][seat.Column]
			seatValid.Taken = false
			seatValid.Group = 0
		}
	}

	return nil
}

// helpers

func (c *cinemaServiceUsecase) validDistance(row, col int, bookingSet map[[2]int]bool) bool {
	for i := 0; i < c.cinema.Rows; i++ {
		for j := 0; j < c.cinema.Cols; j++ {

			if bookingSet[[2]int{i, j}] {
				continue
			}

			lock := c.cinema.SeatLocks[i][j]
			lock.RLock()
			seat := c.cinema.Seats[i][j]
			taken := seat.Taken
			lock.RUnlock()
			if taken && c.manhattan(seat.Row, seat.Column, row, col) < c.cinema.MinDistance {
				return false
			}
		}
	}

	return true
}
func (c *cinemaServiceUsecase) manhattan(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
