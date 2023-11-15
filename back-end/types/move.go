package types

// "encoding/json"

type Move struct {
	Player   int   `json:"player"`
	Reset    bool  `json:"reset"`
	GetMoves bool  `json:"getMoves"`
	To       Point `json:"to"`
	From     Point `json:"from"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPoint(x int, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}
