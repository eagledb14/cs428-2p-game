package types

type BoardUpdate struct {
  ValidMove bool `json:"validMove"`
  PlayerMoveId int `json:"playerMoveId"`
  PlayerTurn int `json:"playerTurn"`
  IsOver bool `json:"isOver"`
  Board []int `json:"board"`
}

func NewBoardUpdate(validMove bool, playerId int, board Board, playerTurn int) BoardUpdate {
  return BoardUpdate {
    ValidMove: validMove,
    PlayerMoveId: playerId,
    PlayerTurn: playerTurn,
    IsOver: false,
    Board: board.Board,
  }
}

func (b *BoardUpdate) Reset(validMove bool, playerId int, board Board) {
  b.ValidMove = validMove
  b.PlayerMoveId = playerId
  b.Board = board.Board
}
