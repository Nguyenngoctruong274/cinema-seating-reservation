package cinema_service

import (
	"cinema-seat-reservation/service/model/request"
	"sync"
)

type CinemaService struct {
	Rows         int
	Cols         int
	MinDistance  int
	Seats        [][]*request.Seat
	SeatLocks    [][]*sync.Mutex
	GroupCounter int
	Mu           sync.Mutex
}

func NewCinemaService() *CinemaService {
	return &CinemaService{}
}

func (s *CinemaService) Configure(rows, cols, minDist int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Rows, s.Cols, s.MinDistance = rows, cols, minDist
	s.Seats = make([][]*request.Seat, rows)
	s.SeatLocks = make([][]*sync.Mutex, rows)
	for i := range s.Seats {
		s.Seats[i] = make([]*request.Seat, cols)
		s.SeatLocks[i] = make([]*sync.Mutex, cols)
		for j := 0; j < cols; j++ {
			s.Seats[i][j] = &request.Seat{Row: i, Column: j}
			s.SeatLocks[i][j] = &sync.Mutex{}
		}
	}
	s.GroupCounter = 0
}
