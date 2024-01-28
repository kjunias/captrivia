import { useContext, useState } from "react";
import { API_BASE } from "../config";
import { GameContext } from "../contexts/GameContext";

const Home = () => {
  const { error, setError, loading, setLoading, setGameSession } =
    useContext(GameContext);

  const startGame = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const data = await res.json();
      setGameSession(data.sessionId);
    } catch (err) {
      setError("Failed to start game.");
    }
    setLoading(false);
  };

  if (error) return <div className="error">Error: {error}</div>;
  if (loading) return <div className="loading">Loading...</div>;

  return <button onClick={startGame}>Start Game</button>;
};
export default Home;
