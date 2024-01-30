import { useContext } from "react";
import { GameContext } from "../contexts/GameContext";
import useGame from "../hooks/useGame";

const Home = () => {
  const { error, loading } =
    useContext(GameContext);
  const {startGame, createRoom} = useGame();

  if (error) return <div className="error">Error: {error}</div>;
  if (loading) return <div className="loading">Loading...</div>;

  return (
    <>
      <button onClick={startGame}>Start Game</button>
      <button onClick={createRoom}>Create Game Room</button>
    </>
  );
};
export default Home;
