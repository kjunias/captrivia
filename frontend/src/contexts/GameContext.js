import { createContext, useState } from "react";

export const GameContext = createContext({});

export const GameProvider = ({ children }) => {
  const [playerID, setPlayerID] = useState(null);
  const [gameSession, setGameSession] = useState(null);
  const [gameRoom, setGameRoom] = useState(null);
  const [updating, setUpdating] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  return (
    <GameContext.Provider
      value={{
        playerID,
        setPlayerID,
        gameSession,
        setGameSession,
        gameRoom,
        setGameRoom,
        updating,
        setUpdating,
        loading,
        setLoading,
        error,
        setError,
      }}
    >
      {children}
    </GameContext.Provider>
  );
};
