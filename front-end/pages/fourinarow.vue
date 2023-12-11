<template>
  <div class="home">
    <h1>{{ player === turn ? "Your turn" : "Waiting for opponent" }}</h1>
    <h1>Lobby ID: {{ this.lobbyId }}</h1>
    <div ref="board" class="connect-four-board">
      <div v-for="(column, colIndex) in table" class="column">
        <div v-for="(cell, rowIndex) in column" :class="cellClass(cell)" @click="dropDisk(colIndex)" class="cell">
          <!-- The img tag should be inside this div -->
          <img v-if="cell !== 0" :src="getImage(cell)" alt="Game piece" class="game-piece"/>
        </div>
      </div>
    </div>
    <div v-if="isOver" class="game-over">
      <h1>{{ winnerMessage }}</h1>
    </div>
    <div class="buttons">
      <button @click="restart()">Restart</button>
      <button @click="copyGameLink()">Share Game</button>
    </div>
      <div class="scores">
          <!-- Red Player Score -->
          <div class="score" :class="turn === 1 ? 'current-player': ''">
              <img :src="redPiece" class="display-cell-1">
              <span>: {{ this.score1 }}</span>
          </div>
          <div class="score tie-score">
              <span>Ties: {{ this.ties }}</span>
          </div>
          <div class="score" :class="turn === 2 ? 'current-player': ''">
            <img :src="yellowPiece" class="display-cell-2">
              <span>: {{ this.score2 }}</span>
          </div>
      </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      table: Array(7).fill().map(() => Array(6).fill(0)), // 7 columns, 6 rows
      turn: 1, // Player 1 starts
      isOver: false,
      socket: undefined,
      winner: 0,
      score1: 0,
      score2: 0,
      ties: 0,
      lobbyId: 0,
      player: 0,
      game: 'fourinarow',
      api: 'game.blackman.zip/api',
      redPiece: '/Red Checker.svg',
      yellowPiece: '/Yellow Checker.svg',
      // Other necessary data
    }
  },
  async mounted() {
    if (this.$route.query.lobbyId) {
      await fetch(`https://${this.api}/${this.game}?lobbyId=${this.$route.query.lobbyId}`)
        .then(response => response.json())
        .then(data => {
            console.log(data)
            this.lobbyId = this.$route.query.lobbyId
            this.table = this.convertBoard(data.board)
        });
    } else {
        await fetch(`https://${this.api}/${this.game}`)
          .then(response => response.json())
          .then(data => {
              console.log(data)
              this.lobbyId = data
              // this.$router.replace({ query: { lobbyId: this.lobbyId} })
          });
    }

    this.setupSocket()

  },
  methods: {
    setupSocket() {
      this.socket = new WebSocket(`wss://${this.api}/ws?lobbyId=${this.lobbyId}`)
      this.socket.onmessage = (event) => {
          const message = JSON.parse(event.data);
          if (!this.player) {
              this.player = message
          } else if (message?.board && message?.validMove) {
              this.table = this.convertBoard(message.board)
              this.turn = message.playerTurn
              this.winner = message.playerMoveId
              this.isOver = message.isOver
          }
          console.log('message', message)
          // Do things based on the event data
      }
      this.socket.onopen = (event) => {
          console.log('opened', event)
          this.connected = true
      }
      this.socket.onclose = async (event) => {
          this.connected = false
          console.log('connection closed')
      }
    },
    convertBoard(board) {
      console.log(board)
      let table = []
      for (let i = 0; i < board.length; i += 7) {
        table.push(board.slice(i, i + 7))
      }
      table = table[0].map((_, colIndex) => table.map(row => row[colIndex]));
      return table;
    },
    dropDisk(colIndex) {
      console.log(colIndex)
      // Game logic to drop a disk into the column
      // Update the table, check for win/tie conditions
      this.socket.send(JSON.stringify({Player: this.player, Reset: false, To: {X: colIndex, Y: colIndex, From: {X: 0, Y: 0}}}))
    },
    restart() {
      this.socket.send(JSON.stringify({ Player: this.player, Reset: true, To: { X: 0, Y: 0 }, From: { X: 0, Y: 0 }}))    
    },
    copyGameLink() {
      // Logic to copy the game link to clipboard
      navigator.clipboard.writeText(`${window.location.href}?lobbyId=${this.lobbyId}`);
    },
    cellClass(cellValue) {
      // Return appropriate class based on cell value
      return {
        'empty-cell': cellValue === 0,
        'player1-cell': cellValue === 1,
        'player2-cell': cellValue === 2,
      };
    },

    getImage(cellValue) {
      if (cellValue === 1) {
        return this.redPiece; // Path to Player 1's image
      } else if (cellValue === 2) {
        return this.yellowPiece; // Path to Player 2's image
      }
      return ''; // Return empty string for empty cells
    },
    // Additional methods as needed
  },
  computed: {
    winnerMessage() {
      if (this.winner === 1) {
        this.score1 += 1;
      } else if (this.winner === 2) {
        this.score2 += 1;
      } else {
        this.ties++;
      }

      return this.winner === 1 ? 'Player 1 wins!' :
          this.winner === 2 ? 'Player 2 wins!' : 'Tie Game!';
    },
    // Other computed properties
  },
  // Lifecycle hooks, etc.
}
</script>
<style>
body {
  font-family: 'Arial', sans-serif; /* or use Google Fonts */
  background-color: #f4f4f4; /* soft background color */
  color: #333; /* readable text color */
}

.connect-four-board {
  display: grid;
  grid-template-columns: repeat(7, min(11vmin, 60px)); /* Increase the size of each column */
  grid-gap: 6px; /* Increase the gap if desired */
  background-color: blue; /* Set the board background to blue */
  padding: 10px; /* Add some padding around the board */
  border-radius: 10px; /* Optional: rounded corners for the board */
  justify-content: center; /* Center the board horizontally */
  margin: auto; /* Also helps in centering the board */
}
.cell {
  height: 11vmin;
  width: 11vmin;
  max-width: 60px;
  max-height: 60px;
  border-radius: 50%;
  background-color: white; /* Change cell color for contrast */
  display: flex;
  justify-content: center;
  align-items: center;
  border: 2px solid #000; /* Optional: add a border to each cell */
}
.cell img {
  max-width: 100%;
  max-height: 100%;
  border-radius: 50%; /* optional, for rounded images */
}

.display-cell-1 {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: red;
  background-size: cover;
  display: flex;
  justify-content: center;
  align-items: center;
}
.display-cell-2 {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: yellow;
  background-size: cover;
  display: flex;
  justify-content: center;
  align-items: center;
}
.player1-cell {
  background-color: red;

}
.player2-cell {
  background-color: yellow;

}
.buttons {
  display: flex;
  justify-content: space-between;
  margin: 10px;
}
button {
    width: 100px;
    cursor: pointer;
    background: rgb(51, 51, 51);
    color: rgb(255, 255, 255);
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: bold;
    min-height: 48px;
    margin-right: 10px;
}
button:last-child {
  margin-right: 0;
}

.buttons button:hover {
  background-color: #0056b3; /* darker shade on hover */
}

/* Additional styling */

.tie-score {
  font-weight: bold;
  color: #333; /* Distinct color */
}

.game-piece {
  max-width: 100%;
  max-height: 100%;
  border-radius: 50%; /* optional, for rounded images */
}
</style>
