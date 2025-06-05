package main

import (
	"cinema-seat-reservation/service/cinema_service"
	"cinema-seat-reservation/service/handler"
	"cinema-seat-reservation/service/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	cinemaService := cinema_service.NewCinemaService()
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

	r.Run(":8080")
}
