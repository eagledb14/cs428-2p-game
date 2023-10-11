<template>
    <div class="home">
        <h1>Welcome to tictactoe!</h1>
        <div class="tic-tac-toe-board">
            <table>
                <tbody>
                    <tr v-for="(row, rowIndex) in table">
                        <td v-for="(item, index) in row" @click="selectedItem(rowIndex, index)">
                            <svg-icon v-if="item === 1" type="mdi" :path="xIcon"></svg-icon>
                            <svg-icon v-if="item === 0" type="mdi" :path="oIcon"></svg-icon>
                        </td>
                    </tr>
                </tbody>
            </table>
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
            table: [[-1, -1, -1],
                    [-1, -1, -1],
                    [-1, -1, -1]],
            xIcon: mdiWindowClose,
            oIcon: mdiCircleOutline,
            prevTable: [[-1, -1, -1],
                        [-1, -1, -1],
                        [-1, -1, -1]],
            turn: 0,
            socket: undefined,
            connected: false,
        }
    },
    mounted() {
        // check for query parameter to join websocket
        // create lobby otherwise, show share link/popup info
        this.socket = new WebSocket('ws://localhost:3001')
        this.socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log('message', message)
            // Do things based on the event data
        }
        this.socket.onopen = (event) => {
            console.log('opened', event)
            this.connected = true
            this.socket.send(JSON.stringify({message: 'sending to server'}))
        }
        this.socket.onclose = (event) => {
            this.connected = false
            console.log('connection closed')
        }
    },
    methods: {
        selectedItem(row, column) {
            this.prevTable = this.table
            if (this.table[row][column] === -1) {
                this.table[row][column] = this.turn
                this.turn = this.turn ? 0 : 1
            } else {
                alert('Invalid Move, Try again!')
            }
        },
        restart() {
            this.table = [[-1, -1, -1],
                            [-1, -1, -1],
                            [-1, -1, -1]]
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