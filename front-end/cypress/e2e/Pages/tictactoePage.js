class tictactoePage {
    Page_elements = {
        mainPage : () => cy.get(`div[class="home"]`),
        heading1 : ()=> cy.get(`h1[data-v-inspector="pages/tictactoe.vue:3:9"]`),
        gameBoard:() => cy.get(`div[class="tic-tac-toe-board"]`),
        scoreBoard : ()=> cy.get(`div[data-v-inspector="pages/tictactoe.vue:27:9"]`),
        lobbyID: () => cy.get(`h1[data-v-inspector="pages/tictactoe.vue:32:9"]`)
    }
    Buttons = {
        restartButton:() => cy.get(`button[data-v-inspector="pages/tictactoe.vue:24:9"]`),
        shareGameButton:() => cy.get(`button[data-v-inspector="pages/tictactoe.vue:25:9"]`),
    }

}
export default new tictactoePage();