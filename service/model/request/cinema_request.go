package request

type Seat struct {
	Row    int    `json:"row" required:"true" validate:"required"`
	Column int    `json:"column" required:"true" validate:"required"`
	Taken  bool   `json:"taken"`
	Group  uint64 `json:"group"`
}

type ConfigRequest struct {
	Rows        int `json:"rows" required:"true" validate:"required"`
	Columns     int `json:"columns"  required:"true" validate:"required"`
	MinDistance int `json:"minDistance"  required:"true" validate:"required"`
}

type QueryAvailableSeatRequest struct {
	Count int `json:"count" required:"true"`
}

type CheckAvailableSeatRequest struct {
	Seats []Seat `json:"seats" validate:"required"`
}

type ReserveSeatRequest struct {
	Seats []Seat `json:"seats" validate:"required"`
}

type CancelSeatRequest struct {
	Seats []Seat `json:"seats" validate:"required"`
}
