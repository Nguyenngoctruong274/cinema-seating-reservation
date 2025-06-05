package usecase

import (
	"github.com/stretchr/testify/assert"
	"home-assignment/home-assignments/service/cinema_service"
	"home-assignment/home-assignments/service/model/request"
	"sync"
	"testing"
)

func TestCinemaService_ConcurrentReservations(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)
	req := request.ConfigRequest{
		Rows:        1,
		Columns:     10,
		MinDistance: 2,
	}
	svc.Configure(req)

	coords := []request.Seat{{Row: 0, Column: 0}, {Row: 0, Column: 1}}

	var wg sync.WaitGroup
	numGoroutines := 10
	chanErrs := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			chanErrs <- svc.ReserveSeats(request.ReserveSeatRequest{
				Seats: coords,
			})
		}(i)
	}

	go func() {
		wg.Wait()
		close(chanErrs)
	}()

	successCount := 0
	for err := range chanErrs {
		if err == nil {
			successCount++
		}
	}

	// Only one reservation should succeed
	assert.Equal(t, 1, successCount, "Only one reservation should succeed")
}

type ReserveJob struct {
	Seats  request.ReserveSeatRequest
	Result chan error
}

func TestCinemaService_HighConcurrencyWithWorkerPool(t *testing.T) {
	cinemaService := cinema_service.NewCinemaService()
	svc := NewCinemaUsecase(cinemaService)
	req := request.ConfigRequest{
		Rows:        1,
		Columns:     1000,
		MinDistance: 2,
	}
	svc.Configure(req)

	jobChan := make(chan ReserveJob, 2000)
	workerCount := 10

	// Launch worker pool
	for i := 0; i < workerCount; i++ {
		go func() {
			for job := range jobChan {
				err := svc.ReserveSeats(job.Seats)
				job.Result <- err
			}
		}()
	}

	numRequests := 1000

	coords := request.ReserveSeatRequest{
		Seats: []request.Seat{
			{Row: 0, Column: 0}, {Row: 0, Column: 1},
		},
	}

	successCount := 0
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := make(chan error, 1)
			jobChan <- ReserveJob{Seats: coords, Result: res}
			if err := <-res; err == nil {
				successCount++
			}
		}()
	}

	wg.Wait()
	close(jobChan)

	assert.Equal(t, 1, successCount, "Only one reservation should succeed out of 1000 concurrent attempts")
}
