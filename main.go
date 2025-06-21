package main

import (
	"cinema-seat-reservation/service/cinema_service"
	"cinema-seat-reservation/service/handler"
	"cinema-seat-reservation/service/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const stateFile = "data_cinema.json"

func main() {
	r := gin.Default()
	cinemaService := cinema_service.NewCinemaService()
	// Load state on startup
	if err := cinemaService.LoadSaveDataCinema(stateFile); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to load state: %v", err)
	}

	// Wrap with panic recovery
	defer func() {
		if re := recover(); re != nil {
			log.Printf("Recovered from panic: %v", re)
			_ = cinemaService.SaveDataCinema(stateFile)
		}
	}()

	cinemaUscase := usecase.NewCinemaUsecase(cinemaService)
	h := handler.NewCinemaHandler(cinemaUscase)

	//Configure Cinema Layout
	r.POST("/configure", h.Configure)
	//Query Available Seats
	r.GET("/available-seats", h.QueryAvailableSeats)
	//Check Available Seats
	r.POST("/check-available", h.CheckAvailableSeats)
	//Reserve Seats
	r.POST("/reserve", h.ReserveSeats)
	//Cancel Reservation
	r.POST("/cancel", h.CancelSeats)
	// Handle shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go handleShutdown(cinemaService, stateFile, quit)
	r.Run(":8080")
}

func handleShutdown(cinemaService *cinema_service.CinemaService,
	stateFile string,
	quit <-chan os.Signal) {
	<-quit
	if err := cinemaService.SaveDataCinema(stateFile); err != nil {
		log.Printf("Failed to save state: %v", err)
	}
	os.Exit(0)
}
