package test_example

import (
	"fmt"
	"testing"
)

type stateFn func() stateFn

type LudoGame struct {
	board       [15]int
	players     []*LudoPlayer
	currentTurn int
}

type LudoPlayer struct {
	pieces    []*LudoPiece ///  In the NewLudoGame function, you should initialize each player's pieces slice with the appropriate length.
	color     string
	canMove   bool
	hasWon    bool
	startTile int
}

type LudoPiece struct {
	tile   int
	player *LudoPlayer
}

func NewLudoGame() *LudoGame {
	game := &LudoGame{}
	/// In the game of Ludo, players start their pieces from a designated start tile
	// on the board. The start tile is usually marked with a colored circle or other
	// symbol to indicate which player's pieces start from there.
	/// When a player's turn begins, they roll a die or spin a spinner to determine(他们掷骰子或旋转转轮来确定)
	// how many spaces to move their piece(s) from the start tile onto the main part of the board(棋盘).
	game.players = []*LudoPlayer{
		&LudoPlayer{
			color:     "red",
			startTile: 0,
			pieces:    []*LudoPiece{{tile: 0}}, //pieces:    make([]*LudoPiece, 4),
		},
		&LudoPlayer{
			color:     "blue",
			startTile: 13,
			pieces:    []*LudoPiece{{tile: 13}}, //pieces:    make([]*LudoPiece, 4),

		},
		&LudoPlayer{
			color:     "yellow",
			startTile: 26,
			pieces:    []*LudoPiece{{tile: 26}}, //pieces:    make([]*LudoPiece, 4),
		},
		&LudoPlayer{
			color:     "green",
			startTile: 39,
			pieces:    []*LudoPiece{{tile: 39}}, //pieces:    make([]*LudoPiece, 4),
		},
	}
	game.currentTurn = 0
	for i := 0; i < len(game.board); i++ {
		game.board[i] = -1
	}
	return game
}

func (game *LudoGame) StartState() stateFn {
	fmt.Println("Starting game")
	game.players[0].canMove = true
	return game.NextState
}

func (game *LudoGame) NextState() stateFn {
	player := game.players[game.currentTurn]
	if player.hasWon {
		return game.EndState
	}
	if !player.canMove {
		game.currentTurn = (game.currentTurn + 1) % len(game.players)
		game.players[game.currentTurn].canMove = true
	}
	fmt.Printf("%s's turn\n", player.color)
	return game.WaitForInputState
}

func (game *LudoGame) WaitForInputState() stateFn {
	// wait for user input here
	return game.MovePieceState
}

func (game *LudoGame) MovePieceState() stateFn {
	player := game.players[game.currentTurn]
	piece := player.pieces[0]
	piece.tile = (piece.tile + 1) % len(game.board)
	game.board[piece.tile] = game.currentTurn
	fmt.Printf("%s moved to tile %d\n", player.color, piece.tile)
	if piece.tile == player.startTile {
		player.hasWon = true
		player.canMove = false
	} else {
		player.canMove = false
	}
	return game.NextState
}

func (game *LudoGame) EndState() stateFn {
	fmt.Println("Game over")
	return nil
}

func (game *LudoGame) Run() {
	state := game.StartState
	for state != nil {
		state = state()
	}
}

func TestLudoGame(t *testing.T) {
	game := NewLudoGame()
	game.Run()
}
