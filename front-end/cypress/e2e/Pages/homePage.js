class homePage {
    Page_elements = {
        mainPage : () => cy.get(`div[class="home"]`),
        heading1 : ()=> cy.get(`h1[data-v-inspector="pages/index.vue:3:9"]`),
        gameOption:() => cy.get(`div[class="game-options"]`),
        heading2 : ()=> cy.get(`h1[data-v-inspector="pages/index.vue:10:9"]`),
        lobby: () => cy.get(`*[data-v-inspector="pages/index.vue:11:9"]`)
    }
    Buttons = {
        checkerGameButton:() => cy.get(`*[href="/checkers"]`),
        tictactoeGameButton:() => cy.get(`*[href="/tictactoe"]`),
        submitButton: () => cy.get(`button[class="submit-button"]`)
    }
    InputBoxes = {
        lobbyInput: () => cy.get(`input[name="lobbyID"]`)
    }
}
export default new homePage();