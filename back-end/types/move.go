package types


type Move struct {
  Player int
  To Point
  From Point
}


type Point struct {
  X int
  Y int
}

func NewPoint(x int, y int) Point {
  return Point {
    X: x,
    Y: y,
  }
}




