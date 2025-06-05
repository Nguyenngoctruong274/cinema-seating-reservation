package handler

import (
	"cinema-seat-reservation/service/model/request"
	"cinema-seat-reservation/service/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type CinemaHandler interface {
	Configure(c *gin.Context)
	QueryAvailableSeats(c *gin.Context)
	CheckAvailableSeats(c *gin.Context)
	ReserveSeats(c *gin.Context)
	CancelSeats(c *gin.Context)
}

type cinemaHandler struct {
	CinemaUseCase usecase.CinemaServiceUseCase
}

func NewCinemaHandler(cinemaUsecase usecase.CinemaServiceUseCase) CinemaHandler {
	return &cinemaHandler{CinemaUseCase: cinemaUsecase}
}

func (h *cinemaHandler) Configure(c *gin.Context) {
	req := request.ConfigRequest{}
	validate := validator.New()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	h.CinemaUseCase.Configure(req)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (h *cinemaHandler) QueryAvailableSeats(c *gin.Context) {
	count := c.Query("count")

	seat, _ := strconv.Atoi(count)
	if seat == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request cannot be empty"})
		return
	}

	data := h.CinemaUseCase.QueryAvailableSeats(request.QueryAvailableSeatRequest{Count: seat})

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *cinemaHandler) CheckAvailableSeats(c *gin.Context) {
	req := request.CheckAvailableSeatRequest{}
	validate := validator.New()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	data := h.CinemaUseCase.CheckAvailableSeats(req)

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *cinemaHandler) ReserveSeats(c *gin.Context) {
	req := request.ReserveSeatRequest{}
	validate := validator.New()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	if err := h.CinemaUseCase.ReserveSeats(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (h *cinemaHandler) CancelSeats(c *gin.Context) {
	req := request.CancelSeatRequest{}
	validate := validator.New()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	if err := h.CinemaUseCase.CancelSeats(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
