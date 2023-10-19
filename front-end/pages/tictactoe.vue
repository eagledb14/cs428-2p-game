<template>
    <div class="home">
        <h1 v-if="!player">Welcome to tictactoe!</h1>
        <h1 v-if="player === turn">Your turn</h1>
        <h1 v-if="player && player !== turn">Waiting for opponent</h1>
        <div class="tic-tac-toe-board">
            <table>
                <tbody>
                    <tr v-for="(row, rowIndex) in table">
                        <td v-for="(item, index) in row" @click="selectedItem(rowIndex, index)">
                            <svg-icon v-if="item === 1" type="mdi" :path="xIcon"></svg-icon>
                            <svg-icon v-if="item === 2" type="mdi" :path="oIcon"></svg-icon>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div v-if="isOver" class="game-over">
            <h1 v-if="player === winner">You win!</h1>
            <h1 v-else-if="winner === -1">Tie Game!</h1>
            <h1 v-else>You loose!</h1>
        </div>
        <button @click="restart()">restart</button>
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
            table: [[0, 0, 0],
                    [0, 0, 0],
                    [0, 0, 0]],
            xIcon: mdiWindowClose,
            oIcon: mdiCircleOutline,
            prevTable: [[0, 0, 0],
                        [0, 0, 0],
                        [0, 0, 0]],
            turn: 1,
            socket: undefined,
            connected: false,
            player: 0,
            isOver: false,
            winner: 0,
        }
    },
    mounted() {
        // check for query parameter to join websocket
        // create lobby otherwise, show share link/popup info
        // We'll need to call something here to get the board/lobby once that is set up
        // fetch('http://localhost:8080/tictactoe')
        //     .then(response => response.json())
        //     .then(data => console.log(data));
        this.socket = new WebSocket('ws://localhost:8080/ws?lobbyId=12345')
        this.socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            if (!this.player) {
                this.player = message
            } else if (message?.board && message?.validMove) {
                // handle move event
                // ## BoardUpdate
                // - params
                // - ValidMove: bool
                // - Opponent: int
                // - isOver: bool 
                // - Board: []int
                this.table = this.convertBoard(message.board)
                if (!message.isOver && !this.isOver) {
                    this.turn = message.playerMoveId === 1 ? 2 : 1
                } 
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
        this.socket.onclose = (event) => {
            this.connected = false
            console.log('connection closed')
        }
    },
    methods: {
        selectedItem(row, column) {
            this.prevTable = this.table
            if (this.table[row][column] === 0 && this.turn === this.player) {
                // ## Move
                // - params
                // - Player: int
                // - Reset: bool
                // - To: Point
                // - From: Point
                this.socket.send(JSON.stringify({ Player: this.player, Reset: false, To: { X: row, Y: column }, From: { X: row, Y: column }}))
            }
        },
        restart() {
            this.table = [[0, 0, 0],
                        [0, 0, 0],
                        [0, 0, 0]]
            this.socket.send(JSON.stringify({ Player: this.player, Reset: true, To: { X: 0, Y: 0 }, From: { X: 0, Y: 0 }}))    
        },
        convertBoard(board) {
            const table = []
            table.push(board.slice(0, 3))
            table.push(board.slice(3, 6))
            table.push(board.slice(6, 9))
            return table
        }
    }
}
</script>
<style>
td {
    height: 55px;
    width: 55px;
    padding: 8px;
    position: relative;
    vertical-align: top;
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
}
td svg {
    height: 92%;
    width: 100%;
}
table {
  border-collapse: collapse;
  border-style: hidden;
}
table td {
  border: 5px solid black;
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
    padding: 48px;
}
</style>