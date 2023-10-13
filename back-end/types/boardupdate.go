package types

type BoardUpdate struct {
  ValidMove bool
  PlayerIdOfMove int
  IsOver bool
  Board []int
}

func NewBoardUpdate(validMove bool, playerId int, board Board) BoardUpdate {
  return BoardUpdate {
    ValidMove: validMove,
    PlayerIdOfMove: playerId,
    IsOver: false,
    Board: board.board,
  }
}

func (b *BoardUpdate) Reset(validMove bool, playerId int, board Board) {
  b.ValidMove = validMove
  b.PlayerIdOfMove = playerId
  b.Board = board.board
}
