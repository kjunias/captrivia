import { useContext, useEffect } from "react";
import { GameContext } from "../contexts/GameContext";
import { API_BASE} from "../config";

const useGame = () => {
  const {
    setError,
    loading,
    setLoading,
    setGameSession,
    setGameRoom,
    setPlayerID
  } = useContext(GameContext);

  useEffect(()=>{
    const controller = new AbortController();
    const joinRoom = async() => {
      if(window.location.pathname == "/join") {
        setLoading(true);
        setError(null);
        try {
          const roomID = window.location.search.split("roomID=")[1].substring(0, 6)
          const res = await fetch(`${API_BASE}/gameroom/join?roomID=${roomID}`, {
            signal: controller.signal,
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
          });
          const data = await res.json();
          setGameRoom(data);
          setPlayerID(data.playerID);
        } catch (err) {
          setError("Failed to join game room");
        } finally {
          setLoading(false);
        }  
      }
    }
    if (!loading) {
      joinRoom();
    }
    return () => controller.abort()
  }, []);

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
      setGameRoom({...data, currentQuestionIndex: -1});
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