<template>
    <div class="home" draggable="false">
        <h1 v-if="!player || !connected">Share Game to Play</h1>
        <h1 v-else-if="player === turn">Your turn</h1>
        <h1 v-else-if="player && player !== turn">Waiting for opponent</h1>
        <div ref="board" class="checker-board" draggable="false">
            <div v-for="(row, rowIndex) in table" class="board-row">
                <div :ref="`${rowIndex},${colIndex}`" v-for="(item, colIndex) in row" @click="selectItem(rowIndex, colIndex)" :class="item === -1 ? 'red-square' : 'black-square'" class="square" draggable="false" :rowIndex="rowIndex" :colIndex="colIndex" >
                    <img v-if="item === 1" src="/Red Checker.svg" class="checker" draggable="false"/>
                    <img v-if="item === 2" src="/Black Checker.svg" class="checker" draggable="false"/>
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
            <div class="score">
                <svg-icon type="mdi" :path="xIcon"></svg-icon>
                <span>: {{ this.score1 }}</span>
            </div>
            <div class="score">
                <span>Ties: {{ this.ties }}</span>
            </div>
            <div class="score">
                <svg-icon type="mdi" :path="oIcon"></svg-icon>
                <span>: {{ this.score2 }}</span>
            </div>
        </div>
        <h1>Lobby ID: {{ this.lobbyId }}</h1>
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
            table: [[-1, 1, -1, 1, -1, 1, -1, 1],
                    [1, -1, 1, -1, 1, -1, 1, -1],
                    [-1, 1, -1, 1, -1, 1, -1, 1],
                    [0, -1, 0, -1, 0, -1, 0, -1],
                    [-1, 0, -1, 0, -1, 0, -1, 0],
                    [2, -1, 2, -1, 2, -1, 2, -1],
                    [-1, 2, -1, 2, -1, 2, -1, 2],
                    [2, -1, 2, -1, 2, -1, 2, -1]],
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
            lobbyId: 0,
            score1: 0,
            score2: 0,
            ties: 0,
            api: 'game.blackman.zip/api',
            game: 'checkers',
            selectedItem: undefined // { row: 0, column: 0 }
        }
    },
    async mounted() {
        // check for query parameter to join websocket
        // create lobby otherwise, show share link/popup info
        // We'll need to call something here to get the board/lobby once that is set up
        // if (this.$route.query.lobbyId) {
        //     await fetch(`https://${this.api}/${this.game}?lobbyId=${this.$route.query.lobbyId}`)
        //         .then(response => response.json())
        //         .then(data => {
        //             console.log(data)
        //             this.lobbyId = this.$route.query.lobbyId
        //             this.table = this.convertBoard(data.board)
        //         });
        // } else {
        //     await fetch(`https://${this.api}/${this.game}`)
        //         .then(response => response.json())
        //         .then(data => {
        //             console.log(data)
        //             this.lobbyId = data
        //             // this.$router.replace({ query: { lobbyId: this.lobbyId} })
        //         });
        // }

        // this.setupSocket()
        this.setupEventListeners()
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
        setupEventListeners() {
            const pointerMove = (event) => {
                const box = this.$refs[`${this.selectedItem.row},${this.selectedItem.col}`][0].getBoundingClientRect()
                const x = event.x - box.x - box.width / 2
                const y = event.y - box.y - box.height / 2
                if (this.$refs[`${this.selectedItem.row},${this.selectedItem.col}`][0].children[0]) {
                    this.$refs[`${this.selectedItem.row},${this.selectedItem.col}`][0].children[0].style = `position: absolute; top: ${y}px; left: ${x}px; z-index: 1000;`
                } else {
                    this.$refs['board'].removeEventListener('pointermove', pointerMove)
                }
            }
            for (let row = 0; row < this.table.length; row++) {
                for (let col = 0; col < this.table[row].length; col++) {
                    const pointerDown = () => {
                        // if (this.table[row][col] === this.player && this.turn === this.player) {
                        if (this.table[row][col] > 0 && !this.selectedItem) {
                            this.selectedItem = { row, col, piece: this.table[row][col] }
                            // request valid moves
                            this.$refs['board'].addEventListener('pointermove', pointerMove)
                        }
                    }
                    this.$refs[`${row},${col}`][0].addEventListener('pointerup', (event) => {
                        this.$refs['board'].removeEventListener('pointermove', pointerMove)
                        // if (this.table[row][col] === this.player && this.turn === this.player) {
                        if (this.selectedItem) {
                            this.$refs[`${this.selectedItem.row},${this.selectedItem.col}`][0].children[0].style = ''
                            const selectedEl = document.elementFromPoint(event.clientX, event.clientY)
                            const selectedRow = selectedEl.getAttribute("rowindex")
                            const selectedCol = selectedEl.getAttribute("colindex")
                            // check if valid move
                            this.$nextTick(() => {
                                if (this.table[selectedRow][selectedCol] === 0) {
                                    // Send move to server
                                    this.table[this.selectedItem.row][this.selectedItem.col] = 0
                                    this.table[selectedRow][selectedCol] = this.selectedItem.piece
                                }
                                this.selectedItem = undefined
                            })
                        }
                    })
                    this.$refs[`${row},${col}`][0].addEventListener('pointerdown', pointerDown)
                }
            }
        },
        selectItem(row, column) {
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
            this.socket.send(JSON.stringify({ Player: this.player, Reset: true, To: { X: 0, Y: 0 }, From: { X: 0, Y: 0 }}))    
        },
        convertBoard(board) {
            const table = []
            table.push(board.slice(0, 3))
            table.push(board.slice(3, 6))
            table.push(board.slice(6, 9))
            return table
        },
        copyGameLink() {
            // navigator.clipboard.writeText(window.location.href);
            navigator.clipboard.writeText(`${window.location.href}?lobbyId=${this.lobbyId}`);
        }
    }
}
</script>
<style>
.square {
    height: 10vmin;
    width: 10vmin;
    max-width: 60px;
    max-height: 60px;
    padding: 1vmin;
    position: relative;
    vertical-align: top;
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
}
.board-row {
    display: flex;

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
.checker-board {
    margin: 0;
    border: 5px solid black;
    touch-action: none;
    user-select: none;
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
.buttons {
    display: flex;
}
.buttons button {
    margin: 10px;
}
.checker {
    width: 100%;
    height: 100%;
    justify-self: center;
    align-self: center;
    pointer-events: none;
    touch-action: none;
    user-select: none;
}
.red-square {
    background-color: red;
}
.black-square {
    background-color: black;
}
</style>