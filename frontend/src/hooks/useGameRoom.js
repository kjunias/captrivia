import { useContext, useEffect, useState } from "react";
import { GameContext } from "../contexts/GameContext";
import { WS_BASE } from "../config";
import useWebSocket, { ReadyState } from 'react-use-websocket';

const COUNT_DOWN_TIME = 5*1000;
const COUNT_DOWN_INTERVAL = 1000;

const useGameRoom = () => {
  const {playerID, gameRoom, setGameRoom, updating, setUpdating} = useContext(GameContext);
  const [numberOfQuestions, setNumberOfQuestions] = useState(0);
  const [timeLeft, setTimeLeft] = useState(COUNT_DOWN_TIME);
  const [counterOn, setCounterOn] = useState(false);
  const [socketUrl, setSocketUrl] = useState(`${WS_BASE}/gameroom/websocket`);
  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);


  const startCountDown = () => {
    setCounterOn(true);
    const timer = setInterval(() => {
      if(timeLeft >= 0) {
        clearInterval(timer);
        setCounterOn(false);
        endCountDown();

      }
      setTimeLeft(Math.max(0, timeLeft - COUNT_DOWN_INTERVAL));
    }, 1000)
  };

  const endCountDown = () => {
    setUpdating(true)
    sendMessage(JSON.stringify({action: "END_COUNTER", roomID: gameRoom.roomID, numberOfQuestions}));
    setUpdating(false)
  };

  const submitAnswer = (index) => {
    const currentQuestion = gameRoom.questions[gameRoom.currentQuestionIndex];
    sendMessage(JSON.stringify({
      action: "SUBMIT_ANSWER",
      roomID: gameRoom.roomID,
      playerID,
      answer: index,
      questionId: currentQuestion.id,
    }));
  }

  useEffect(() => {
    if (gameRoom.isCountingDown && !counterOn) {
      startCountDown();
    }
  },[gameRoom.isCountingDown, counterOn]);
  
  const startGame = async () => {
    setUpdating(true);
    sendMessage(JSON.stringify({
      action: "START_COUNTER", 
      roomID: gameRoom.roomID, 
      playerID,
      numberOfQuestions
    }));
    setUpdating(false);
  };

  useEffect(() => {
    setGameRoom({...gameRoom, currentQuestionIndex: -1})
  }, []);

  useEffect(() => {
    if (lastMessage !== null) {
      const data = JSON.parse(lastMessage.data)
      const updatedData = {...gameRoom, ...data};
      setGameRoom({...gameRoom, ...data})
    }
  }, [lastMessage]);

  return {
    playerID,
    gameRoom,
    setGameRoom,
    startGame,
    numberOfQuestions,
    setNumberOfQuestions,
    timeLeft,
    startCountDown,
    submitAnswer
  };
};

export default useGameRoom;
