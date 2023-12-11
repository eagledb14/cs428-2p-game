<template>
    <div class="home">
        <h1 v-if="!player || !connected">Share Game to Play</h1>
        <h1 v-else-if="player === turn">Your turn</h1>
        <h1 v-else-if="player && player !== turn">Waiting for opponent</h1>
        <h1>Lobby ID: {{ this.lobbyId }}</h1>
        <div class="othello-board">
            <div v-for="(row, rowIndex) in table" class="board-row">
                <div v-for="(item, colIndex) in row" 
                class="square" @click="selectedItem(rowIndex, colIndex)">
                    <img v-if="item === 1" :src="blackPiece" class="checker"/>
                    <img v-if="item === 2" :src="whitePiece" class="checker"/>
                </div>
            </div>
        </div>
        <div v-if="isOver" class="game-over">
            <h1 v-if="player === winner">You win!</h1>
            <h1 v-else-if="winner === -1">Tie Game!</h1>
            <h1 v-else>You lose!</h1>
        </div>
        <div class="buttons">
            <button v-if="isOver" @click="restart()">Play Again</button>
            <button v-else @click="restart()">Restart</button>
            <button @click="copyGameLink()" >Share Game</button>
        </div>

        <div class="scores">
            <div class="score" :class="turn === 1 ? 'current-player': ''">
                <img :src="blackPiece" class="player-icon"/>
                <span>: {{ this.score1 }}</span>
            </div>
            <div class="score">
                <span>Ties: {{ this.ties }}</span>
            </div>
            <div class="score" :class="turn === 2 ? 'current-player': ''">
                <img :src="whitePiece" class="player-icon"/>
                <span>: {{ this.score2 }}</span>
            </div>
        </div>
    </div>
</template>
<script>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCircleOutline, mdiWindowClose } from '@mdi/js';
import Websocket from 'ws';
export default {
    components: {SvgIcon},
    data() {
        return {
            table: [[0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 1, 2, 0, 0, 0],
                    [0, 0, 0, 2, 1, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0]],
            xIcon: mdiWindowClose,
            oIcon: mdiCircleOutline,
            prevTable: [[0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 1, 2, 0, 0, 0],
                    [0, 0, 0, 2, 1, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0],
                    [0, 0, 0, 0, 0, 0, 0, 0]],
            turn: 1,
            socket: undefined,
            connected: false,
            player: 0,
            isOver: false,
            winner: 0,
            lobbyId: 0,
            score1: 0,
            score2: 0,
            ties: 0,
            api: 'game.blackman.zip/api',
            game: 'othello',
            whitePiece: '/White Checker.svg',
            blackPiece: '/Black Checker.svg'
        }
    },
    async mounted() {
        // check for query parameter to join websocket
        // create lobby otherwise, show share link/popup info
        // We'll need to call something here to get the board/lobby once that is set up
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
    watch: {
        isOver() {
            if(this.isOver) {
                if (this.winner === 1) {
                    this.score1++
                } else if (this.winner === 2) {
                    this.score2++
                } else {
                    this.ties++
                }
            }
        }
    },
    methods: {
        setupSocket() {
            this.socket = new WebSocket(`wss://${this.api}/ws?lobbyId=${this.lobbyId}`)
            this.socket.onmessage = (event) => {
                const message = JSON.parse(event.data);
                if (!this.player) {
                    this.player = message
                } else if (message?.board && message?.validMove) {
                    // handle move event
                    // ## BoardUpdate
                    // - params
                    // - ValidMove: bool
                    // - PlayerMoveId: int
                    // - PlayerTurn: playerTurn,
                    // - isOver: bool 
                    // - Board: []int
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
                // await fetch(`https://${this.api}/tictactoe`)
                // .then(response => response.json())
                // .then(data => {
                //     console.log(data)
                //     this.lobbyId = data
                //     this.$router.replace({ query: { lobbyId: this.lobbyId} })
                //     this.setupSocket()
                // });
            }
        },
        selectedItem(row, column) {
            if (this.table[row][column] === 0 && this.turn === this.player) {
                // ## Move
                // - params
                // - Player: int
                // - Reset: bool
                // - To: Point
                // - From: Point
                this.socket.send(JSON.stringify({ Player: this.player, Reset: false, To: { X: column, Y: row }, From: { X: row, Y: column }}))
            }
        },
        restart() {  
            this.socket.send(JSON.stringify({ Player: this.player, Reset: true, To: { X: 0, Y: 0 }, From: { X: 0, Y: 0 }}))    
        },
        convertBoard(board) {
            console.log("board")
            const table = [];
            for (let i = 0; i < board.length; i += 8) {
                table.push(board.slice(i, i + 8));
            };
            return table
        },
        copyGameLink() {
            // navigator.clipboard.writeText(window.location.href);
            navigator.clipboard.writeText(`${window.location.href}?lobbyId=${this.lobbyId}`);
        }
    }
}
</script>
<style scoped="true">
.square {
    height: 10.5vmin;
    width: 10.5vmin;
    max-width: 60px;
    max-height: 60px;
    padding: 0.5vmin;
    position: relative;
    vertical-align: top;
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: rgb(19, 84, 32);
    border: 0.25vmin solid rgb(13, 58, 22);
}
.board-row {
    display: flex;
}
.othello-board {
    margin: 0;
    border: 5px solid black;
}
table {
    border-collapse: collapse;
    border-style: hidden;
}
.checker {
    width: 100%;
    height: 100%;
    justify-self: center;
    align-self: center;
}

tr {
    display: flex;
}
table td {
  border: 5px solid rgb(13, 58, 22);
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
.tic-tac-toe-board {
    padding: 24px 0;
}
.player-icon {
    width: 25px;
    height: 25px;
}
.buttons {
    display: flex;
}
.buttons button {
    margin: 10px;
}
</style>
