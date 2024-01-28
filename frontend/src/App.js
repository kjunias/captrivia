import React, { useContext } from "react";
import "./App.css";
import Home from "./components/Home";
import { GameContext, GameProvider } from "./contexts/GameContext";
import QuestionContainer from "./components/Question";

function AppContainer() {
  const { gameSession } = useContext(GameContext);
  return (
    <div className="App">{!gameSession ? <Home /> : <QuestionContainer />}</div>
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
