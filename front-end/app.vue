<template>
    <div>
        <NuxtPage />
    </div>
</template>
<style>
.home {
    display: flex;
    align-items: center;
    flex-direction: column;
    position: fixed;
    width: 100%;
    height: 100%;
    overflow: scroll;
    /* background-color: #0cc5f7; */
}
body {
    font-family: arial;
    margin: 0;
}
.scores {
    justify-content: center;
    display: flex;
    width: 100%;
    padding: 20px 0;
    align-items: center;
}
.score {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 5px;
    margin-right: 20px; /* Space between score displays */
    font-size: 1.5em; /* Enlarge the font size */
    border-radius: 15px;
}
.scores .score.current-player {
    box-shadow: inset 0px 0 0px 2px black;
}
</style>