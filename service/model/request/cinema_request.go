package request

type Seat struct {
	Row    int    `json:"row" required:"true"`
	Column int    `json:"column" required:"true"`
	Taken  bool   `json:"taken"`
	Group  uint64 `json:"group"`
}

type ConfigRequest struct {
	Rows        int `json:"rows" required:"true"`
	Columns     int `json:"columns"  required:"true"`
	MinDistance int `json:"minDistance"  required:"true"`
}

type QueryAvailableSeatRequest struct {
	Count int `json:"count" required:"true"`
}

type CheckAvailableSeatRequest struct {
	Seats []Seat `json:"seats"`
}

type ReserveSeatRequest struct {
	Seats []Seat `json:"seats"`
}

type CancelSeatRequest struct {
	Seats []Seat `json:"seats"`
}
