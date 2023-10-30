<template>
    <div class="home">
        <h1>
            Select a game below
        </h1>
        <div class="game-options">
            <!-- <NuxtLink to="/checkers" class="game-icon"><img src="/checkers.jpg"></NuxtLink> -->
            <NuxtLink to="/tictactoe" class="game-icon"><img src="/tictactoe.jpg"></NuxtLink>
        </div>
        <h1> or input a lobby ID</h1>
        <form @submit.prevent="joinLobby">
            <input v-model="lobbyId" name="lobbyID" placeholder="Lobby ID" inputmode="numeric" autocomplete="off" >
            <button class="submit-button" type="submit">Enter</button>
        </form>
    </div>
</template>
<script>

export default {
    data() {
        return {
            lobbyId: undefined,
            api: 'game.blackman.zip/api',
        }
    },
    methods: {
        async joinLobby() {
            console.log('join lobby with id: ', this.lobbyId)
            await fetch(`https://${this.api}/lobby?lobbyId=${this.lobbyId}`)
                .then(response => response.text())
                .then(response => {
                    console.log(response)
                    if (response === 'tictactoe') {
                        console.log('push to tictactoe')
                        navigateTo({
                            path: '/tictactoe',
                            query: {
                                lobbyId: this.lobbyId
                            }
                        })
                    }
                })
        }
    }
}

</script>
<style> 
.game-options {
    display: flex;
}
.game-icon {
    width: 100px;
    height: 100px;
    border-radius: 999px;
}
.game-icon img {
    width: 100%;
    height: 100%;
    border-radius: inherit;
}

form {
    background-color: rgb(255, 255, 255);
    padding: 16px;
    border-radius: 4px;
}
input {
    line-height: 2.375rem;
    text-align: center;
    color: rgb(51, 51, 51);
    font-size: 1rem;
    border-radius: 4px;
    box-sizing: border-box;
    width: 100%;
    margin-bottom: 0.625rem;
}
.submit-button {
    width: 100%;
    cursor: pointer;
    background: rgb(51, 51, 51);
    color: rgb(255, 255, 255);
    border-radius: 4px;
    font-size: 0.875rem;
    font-weight: bold;
    min-height: 48px;
}
</style>