package response

type Seat struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Taken  bool   `json:"-"`
	Group  uint64 `json:"-"`
}

type CheckAvailableSeatResp struct {
	AvailableSeats []Seat `json:"available_seats"`
}

type QueryAvailableSeatResp struct {
	AvailableSeats [][]Seat `json:"available_seats"`
}

type ReserveSeatsResp struct {
	Message string `json:"message"`
}
