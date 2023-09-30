<template>
    <div class="home">
        <h1>
            2 player games! Select a game below to begin, or input a game code
        </h1>
        <NuxtLink to="/checkers" class="game-icon"><img src="/checkers.jpg"></NuxtLink>
    </div>
</template>
<style> 
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
.home {
    display: flex;
    align-items: center;
    flex-direction: column;
    position: fixed;
    width: 100%;
    height: 100%;
    background-color: #0cc5f7;
}
</style>