<template>
    <div class="home" draggable="false">
        <h1 v-if="!player || !connected">Share Game to Play</h1>
        <h1 v-else-if="player === turn">Your turn</h1>
        <h1 v-else-if="player && player !== turn">Waiting for opponent</h1>
        <div ref="board" class="checker-board" draggable="false">
            <div v-for="(row, rowIndex) in table" class="board-row">
                <div :ref="`${rowIndex},${colIndex}`" v-for="(item, colIndex) in row" 
                class="square" draggable="false" 
                :class="{ 
                    'red-square' : item === -1, 
                    'black-square': item !== -1, 
                    'valid-move': currentValidMoves[rowIndex][colIndex] === 5 || currentValidMoves[rowIndex][colIndex] === 6}" 
                :rowIndex="rowIndex" 
                :colIndex="colIndex" >
                    <img v-if="item === 2" :src="redPiece" class="checker" draggable="false"/>
                    <img v-if="item === 4" :src="redQueen" class="checker" draggable="false"/>
                    <img v-if="item === 1" :src="blackPiece" class="checker" draggable="false"/>
                    <img v-if="item === 3" :src="blackQueen" class="checker" draggable="false"/>
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
            <button v-if="!connected" @click="copyGameLink()" >Share Game</button>
            <button v-if="isJumpMove" @click="pass()">Pass</button>
        </div>

        <div class="scores">
            <div class="score" :class="turn === 2 ? 'current-player': ''">
                <img :src="redPiece" class="player-icon"/>
                <span>: {{ this.score2 }}</span>
            </div>
            <div class="score">
                <span>Ties: {{ this.ties }}</span>
            </div>
            <div class="score" :class="turn === 1 ? 'current-player': ''">
                <img :src="blackPiece" class="player-icon"/>
                <span>: {{ this.score1 }}</span>
            </div>
        </div>
        <h1>Lobby ID: {{ this.lobbyId }}</h1>
    </div>
</template>
<script>
import SvgIcon from '@jamescoyle/vue-icon';
export default {
    components: {SvgIcon},
    data() {
        return {
            table: [[-1, 2, -1, 0, -1, 1, -1, 1],
                    [2, -1, 2, -1, 0, -1, 1, -1],
                    [-1, 2, -1, 0, -1, 1, -1, 1],
                    [2, -1, 2, -1, 0, -1, 1, -1],
                    [-1, 2, -1, 0, -1, 1, -1, 1],
                    [2, -1, 2, -1, 0, -1, 1, -1],
                    [-1, 2, -1, 0, -1, 1, -1, 1],
                    [2, -1, 2, -1, 0, -1, 1, -1]],
            defaultMoves: [[0, 0, 0, 0, 0, 0, 0, 0],
                           [0, 0, 0, 0, 0, 0, 0, 0],
                           [0, 0, 0, 0, 0, 0, 0, 0],
                           [0, 0, 0, 0, 0, 0, 0, 0],
                           [0, 0, 0, 0, 0, 0, 0, 0],
                           [0, 0, 0, 0, 0, 0, 0, 0],
                           [0, 0, 0, 0, 0, 0, 0, 0],
                           [0, 0, 0, 0, 0, 0, 0, 0]],
            validMoves: undefined,
            previousTurn: -1,
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
            blackPiece: '/Black Checker.svg',
            blackQueen: '/Black Checker Queen.svg',
            redPiece: '/Red Checker.svg',
            redQueen: '/Red Checker Queen.svg',
            selectedItem: undefined, // { row: 0, column: 0 }
            previousSelectedItem: undefined
        }
    },
    async mounted() {
        // check for query parameter to join websocket
        // create lobby otherwise, show share link/popup info
        // We'll need to call something here to get the board/lobby once that is set up
        if (this.$route.query.lobbyId) {
            await fetch(`https://${this.api}/${this.game}?lobbyId=${this.$route.query.lobbyId}`)
                .then(response => {
                    return response.json()
                }).catch(error => {
                    console.error(error)
                    this.$router.replace({ path: this.$route.path })
                })
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
        },
        selectedItem() {
            if (this.selectedItem) {
                // - GetMoves: bool
                const row = this.selectedItem.row
                const column = this.selectedItem.col
                console.log('selectedItem', this.selectedItem)
                const move = { Player: this.player, Reset: false, To: { X: row, Y: column }, From: { X: row, Y: column }, GetMoves: true}
                if (this.isJumpMove) {
                    if (this.previousSelectedItem.row === row && this.previousSelectedItem.column === column){
                        move.JumpOnly = true
                    } else {
                        return
                    }
                }
                console.log(move)
                this.socket.send(JSON.stringify(move))
            }
        }
    },
    computed: {
        currentValidMoves() {
            return this.validMoves || this.defaultMoves
        },
        isJumpMove() {
            return this.turn === this.previousTurn
        }
    },
    methods: {
        setupSocket() {
            this.socket = new WebSocket(`wss://${this.api}/ws?lobbyId=${this.lobbyId}`)
            this.socket.onmessage = (event) => {
                const message = JSON.parse(event.data);
                if (!this.player) {
                    this.player = message
                    this.connected = true
                } else if (message?.board && this.selectedItem) {
                    this.validMoves = this.convertBoard(message.board)                    
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
                    this.previousTurn = this.turn
                    this.turn = message.playerTurn
                    this.winner = message.playerMoveId
                    this.isOver = message.isOver
                }
                console.log('message', message)
                // Do things based on the event data
            }
            this.socket.onopen = (event) => {
                console.log('opened', event)
                // this.connected = true
            }
            this.socket.onclose = async (event) => {
                this.connected = false
                this.$router.replace({ path: this.$route.path })
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
                        if ((this.table[row][col] === this.player || this.table[row][col] === this.player + 2) && !this.selectedItem && this.turn === this.player && !this.isOver) {
                            this.selectedItem = { row, col, piece: this.table[row][col] }
                            // request valid moves
                            this.$refs['board'].addEventListener('pointermove', pointerMove)
                        }
                    }
                    this.$refs[`${row},${col}`][0].addEventListener('pointerup', (event) => {
                        this.$refs['board'].removeEventListener('pointermove', pointerMove)
                        if (this.selectedItem) {
                            this.$refs[`${this.selectedItem.row},${this.selectedItem.col}`][0].children[0].style = ''
                            const selectedEl = document.elementFromPoint(event.clientX, event.clientY)
                            const selectedRow = parseInt(selectedEl.getAttribute("rowindex"))
                            const selectedCol = parseInt(selectedEl.getAttribute("colindex"))
                            // check if valid move
                            this.$nextTick(() => {
                                if (this.table[selectedRow][selectedCol] === 0) {
                                    // Send move to server
                                    this.selectItem({ row: this.selectedItem.row, column: this.selectedItem.col }, { row: selectedRow, column: selectedCol })
                                }
                                this.previousSelectedItem = { row: selectedRow, column: selectedCol }
                                this.selectedItem = undefined
                                this.validMoves = undefined
                            })
                        }
                    })
                    this.$refs[`${row},${col}`][0].addEventListener('pointerdown', pointerDown)
                }
            }
        },
        selectItem(from, to) {
            if (this.table[to.row][to.column] === 0 && this.turn === this.player) {
                // ## Move
                // - params
                // - Player: int
                // - Reset: bool
                // - To: Point
                // - From: Point
                console.log('selectItem', from, to)
                const move = { Player: this.player, Reset: false, To: { X: to.row, Y: to.column }, From: { X: from.row, Y: from.column }, GetMoves: false}
                console.log(move)
                this.socket.send(JSON.stringify(move))
            }
        },
        restart() {  
            this.socket.send(JSON.stringify({ Player: this.player, Reset: true, To: { X: 0, Y: 0 }, From: { X: 0, Y: 0 }}))    
        },
        pass() {
            this.socket.send(JSON.stringify({ Player: this.player, Reset: false, To: { X: 0, Y: 0 }, From: { X: 0, Y: 0 }, Pass: true }))
        },
        convertBoard(board) {
            let table = []
            for (let i = 0; i < board.length; i += 8) {
                table.push(board.slice(i, i + 8))
            }
            table = table[0].map((_, colIndex) => table.map(row => row[colIndex]));
            return table
        },
        copyGameLink() {
            navigator.clipboard.writeText(`${window.location.href}?lobbyId=${this.lobbyId}`);
        }
    }
}
</script>
<style>
.square {
    height: 11vmin;
    width: 11vmin;
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
.valid-move {
    box-shadow: inset 0px 0 0px 0.5vmin yellow;
}
</style>