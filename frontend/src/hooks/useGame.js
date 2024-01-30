import { useContext } from "react";
import { GameContext } from "../contexts/GameContext";
import { API_BASE } from "../config";

const useGame = () => {
  const {setError, setLoading, setGameSession, setGameRoom, setPlayerID} =
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

  const createRoom = async() => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_BASE}/gameroom/create`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const data = await res.json();
      setGameRoom(data.roomID);
      setPlayerID(data.playerID);
    } catch (err) {
      setError("Failed to create game room");
    }
    setLoading(false);
  };
  
  return {
    startGame,
    createRoom,
  };
};

export default useGame;