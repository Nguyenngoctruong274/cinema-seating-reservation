package response

type Seat struct {
	Row    int  `json:"row"`
	Column int  `json:"column"`
	Taken  bool `json:"-"`
	Group  int  `json:"-"`
}

type CheckAvailableSeatResp struct {
	AvailableSeats []Seat `json:"available_seats"`
}

type QueryAvailableSeatResp struct {
	AvailableSeats [][]Seat `json:"available_seats"`
}
