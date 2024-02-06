import React, { useContext } from "react";
import "./App.css";
import Home from "./components/Home";
import { GameContext, GameProvider } from "./contexts/GameContext";
import QuestionContainer from "./components/Questions";
import GameRoom from "./components/GameRoom";

function AppContainer() {
  const { gameSession, gameRoom } = useContext(GameContext);
  let screen = <Home/>;

  if (gameRoom) {
    screen = <GameRoom/>
  } else if (gameSession) {
    screen = <QuestionContainer/>
  }
  
  return (
    <div className="App">{screen}</div>
  );
}

function App() {
  return (
    <GameProvider>
      <AppContainer />
    </GameProvider>
  );
}
export default App;
