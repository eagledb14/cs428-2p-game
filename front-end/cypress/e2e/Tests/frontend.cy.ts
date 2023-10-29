import homePage from "../Pages/homePage";
import tictactoe from "../Pages/tictactoePage";
describe('Front Test', () => {
  beforeEach(() => {
    // Visit a specific URL
    let frontend = Cypress.env("frontend")
    let hostname = frontend.hostname;
    let port= frontend.port;
    cy.visit(`${hostname}${port}`);
    // Get the current URL and assert its value
    cy.url().should('eq', 'http://localhost:3000/');
  });


  it('should visit a specific URL and verify the URL', () => {
    homePage.Page_elements.mainPage().should('exist')
    homePage.Page_elements.heading1().should('contain', ' Select a game below ');
    homePage.Page_elements.gameOption().should('exist')
    homePage.Page_elements.heading2().should('contain', ' or input a lobby ID');
    homePage.Page_elements.lobby().should('exist')
    homePage.InputBoxes.lobbyInput().type("702285")
    homePage.Buttons.submitButton().click()
  });

  it('should visit a specific URL and verify the URL', () => {
    homePage.Buttons.checkerGameButton().click()
    cy.url().should('eq', 'http://localhost:3000/checkers');
  });

  it('should visit a specific URL and verify the URL', () => {
    homePage.Buttons.tictactoeGameButton().click()
    cy.url().should('eq', 'http://localhost:3000/tictactoe');
    tictactoe.Page_elements.mainPage().should('exist')
    tictactoe.Page_elements.gameBoard().should('exist')
    tictactoe.Page_elements.heading1().should('contain', 'Welcome to tictactoe!');
    tictactoe.Page_elements.scoreBoard().should('exist')
    tictactoe.Page_elements.lobbyID().should('contain', 'Lobby ID:');
  });



});
