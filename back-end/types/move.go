package types

import (
  // "encoding/json"
)


type Move struct {
  Player int `json:"player"`
  Reset bool `json:"reset"`
  To Point `json:"to"`
  From Point `json:"from"`
}


type Point struct {
  X int `json:"x"`
  Y int `json:"y"`
}

func NewPoint(x int, y int) Point {
  return Point {
    X: x,
    Y: y,
  }
}

func (p *Point) AddPoint(point Point) {
  p.X += point.X
  p.Y += point.Y
}




