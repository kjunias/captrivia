import { useContext, useEffect, useState } from "react";
import { GameContext } from "../contexts/GameContext";
import { API_BASE } from "../config";

const COUNT_DOWN_TIME = 5*1000;
const COUNT_DOWN_INTERVAL = 1000;

const useGameRoom = () => {
  const { gameRoom, setGameRoom, updating, setUpdating} = useContext(GameContext);
  const [numberOfQuestions, setNumberOfQuestions] = useState(0);
  const [timeLeft, setTimeLeft] = useState(COUNT_DOWN_TIME);
  const [counterOn, setCounterOn] = useState(false);

  const startCountDown = () => {
    setCounterOn(true);
    const timer = setInterval(() => {
      if(timeLeft >= 0) {
        clearInterval(timer);
        setCounterOn(false);

      }
      console.log("Timer:",  timeLeft);
      setTimeLeft(Math.max(0, timeLeft - COUNT_DOWN_INTERVAL));
    }, 1000)
  };

  useEffect(() => {
    console.log("==>isCountingDown:", gameRoom.isCountingDown);
    console.log("==>counterOn:", counterOn);
    if (gameRoom.isCountingDown && !counterOn) {
      console.log("==>startCountDown!!!");
      startCountDown();
    }
  },[gameRoom.isCountingDown, counterOn]);
  
  const startGame = async () => {
    try {
      setUpdating(true)
      console.log("Start game start:");
      const res = await fetch(`${API_BASE}/gameroom/start`,{
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Accept": "application/json",
        },
        body: JSON.stringify({
          roomID: gameRoom.roomID,
          numberOfQuestions,
        }),
      });
      console.log("Start game response:", res);
      const content =  await new Response(res.body).text();
      console.log("Start game content:", content);

      if (res.status < 400 && content) {
        setGameRoom(JSON.parse(content));
      }
      setUpdating(false)
    } catch (err) {
      console.error("Failed to start game: ", "err:", err);
      setUpdating(false)
    }
  };

  useEffect(() => {
      const fetchUpdates = async function () {
        if (updating) {
          console.warn("===> updating...")
          return;
        }
        setUpdating(true)
        try {
          const res = await fetch(`${API_BASE}/gameroom/update?roomID=${gameRoom.roomID}&playerID=${gameRoom.playerID}`,{
            method: "GET",
            headers: {
              "Content-Type": "application/json",
              "Accept": "application/json",
            },
          });
          const content =  await new Response(res.body).text()
          if (res.status < 400 && content) {
            setGameRoom(JSON.parse(content));
          }
          setUpdating(false)
        } catch (err) {
          console.error("Failed to fetch updates: ", "err:", err);
          setUpdating(false)
        } finally {
          fetchUpdates();
        }
      };
      fetchUpdates();
  }, []);

  return {
    gameRoom,
    setGameRoom,
    startGame,
    numberOfQuestions,
    setNumberOfQuestions,
    timeLeft,
    startCountDown
  };
};

export default useGameRoom;
