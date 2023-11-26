<template>
  <div class="home">
    <h1>{{ player === turn ? "Your turn" : "Waiting for opponent" }}</h1>
    <h1>Lobby ID: {{ this.lobbyId }}</h1>
    <div ref="board" class="connect-four-board">
      <div v-for="(column, colIndex) in table" class="column">
        <div v-for="(cell, rowIndex) in column" :class="cellClass(cell)" @click="dropDisk(colIndex)" class="cell"></div>
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
          <div class="score">
              <div class="display-cell-1"></div>
              <span>: {{ this.score1 }}</span>
          </div>
          <div class="score">
              <span>Ties: {{ this.ties }}</span>
          </div>
          <div class="score">
              <div class="display-cell-2" ></div>
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
      const table = [];
      for (let i = 0; i < board.length; i += 7) {
        table.push(board.slice(i, i + 7))
      };
      return table[0].map((_, colIndex) => table.map(row => row[colIndex]));
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
.connect-four-board {
  display: grid;
  grid-template-columns: repeat(7, 50px);
  grid-gap: 5px;
}
.cell {
  width: 5.5vh;
  height: 5.5vh;
  border-radius: 50%;
  background-color: lightgray;
  display: flex;
  justify-content: center;
  align-items: center;
}
.display-cell-1 {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: red;
  display: flex;
  justify-content: center;
  align-items: center;
}
.display-cell-2 {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: yellow;
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
}
.score {
    display: flex;
    justify-content: center;
    align-items: center;
}
.scores {
    justify-content: space-between;
    display: flex;
    width: 150px;
    padding: 20px;
}
/* Additional styling */
</style>
