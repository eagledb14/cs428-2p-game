package types

type BoardUpdate struct {
  ValidMove bool
  PlayerIdOfMove int
  Board []int
}

func NewBoardUpdate(validMove bool, playerId int, board Board) BoardUpdate {
  return BoardUpdate {
    ValidMove: validMove,
    PlayerIdOfMove: playerId,
    Board: board.board,
  }
}

func (b *BoardUpdate) Reset(validMove bool, playerId int, board Board) {
  b.ValidMove = validMove
  b.PlayerIdOfMove = playerId
  b.Board = board.board
}
