import { createContext, useState } from "react";

export const GameContext = createContext({});

export const GameProvider = ({ children }) => {
  const [gameSession, setGameSession] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  return (
    <GameContext.Provider
      value={{
        gameSession,
        setGameSession,
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
