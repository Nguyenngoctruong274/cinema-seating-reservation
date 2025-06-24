package cinema_service

import (
	"cinema-seat-reservation/service/model/request"
	"encoding/json"
	"os"
	"sync"
)

type CinemaService struct {
	Rows         int
	Cols         int
	MinDistance  int
	Seats        [][]*request.Seat
	SeatLocks    [][]*sync.RWMutex
	GroupCounter uint64
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
	s.SeatLocks = make([][]*sync.RWMutex, rows)
	for i := range s.Seats {
		s.Seats[i] = make([]*request.Seat, cols)
		s.SeatLocks[i] = make([]*sync.RWMutex, cols)
		for j := 0; j < cols; j++ {
			s.Seats[i][j] = &request.Seat{Row: i, Column: j}
			s.SeatLocks[i][j] = &sync.RWMutex{}
		}
	}
	s.GroupCounter = 0
}

func (s *CinemaService) SaveDataCinema(filename string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	type State struct {
		Rows         int
		Cols         int
		MinDistance  int
		Seats        [][]*request.Seat
		GroupCounter uint64
	}

	if s.Rows == 0 || s.Cols == 0 {
		return nil // Không có dữ liệu để lưu
	}

	state := State{
		Rows:         s.Rows,
		Cols:         s.Cols,
		MinDistance:  s.MinDistance,
		Seats:        s.Seats,
		GroupCounter: s.GroupCounter,
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *CinemaService) LoadSaveDataCinema(filename string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	type State struct {
		Rows         int
		Cols         int
		MinDistance  int
		Seats        [][]*request.Seat
		GroupCounter uint64
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var state State
	if err = json.Unmarshal(data, &state); err != nil {
		return err
	}
	s.Rows = state.Rows
	s.Cols = state.Cols
	s.MinDistance = state.MinDistance
	s.Seats = state.Seats
	s.GroupCounter = state.GroupCounter
	// Khởi tạo lại SeatLocks
	s.SeatLocks = make([][]*sync.RWMutex, s.Rows)
	for i := 0; i < s.Rows; i++ {
		s.SeatLocks[i] = make([]*sync.RWMutex, s.Cols)
		for j := 0; j < s.Cols; j++ {
			s.SeatLocks[i][j] = &sync.RWMutex{}
		}
	}
	return nil
}
