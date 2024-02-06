import { useContext, useEffect, useState } from "react";
import { GameContext } from "../contexts/GameContext";
import { WS_BASE } from "../config";
import useWebSocket from 'react-use-websocket';
import { GameRoomState } from "../constants";

const COUNT_DOWN_TIME = 5*1000;
const COUNT_DOWN_INTERVAL = 1000;

const useGameRoom = () => {
  const {playerID, gameRoom, setGameRoom, setUpdating} = useContext(GameContext);
  const [numberOfQuestions, setNumberOfQuestions] = useState(0);
  const [timeLeft, setTimeLeft] = useState(COUNT_DOWN_TIME);
  const [counterOn, setCounterOn] = useState(false);
  const [socketUrl, _] = useState(`${WS_BASE}/gameroom/websocket`);
  const { sendMessage, lastMessage } = useWebSocket(socketUrl);


  const startCountDown = () => {
    setCounterOn(true);
    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev === 0) {
          endCountDown();
          setCounterOn(false);
          clearInterval(timer);
        }
        return Math.max(0, prev - COUNT_DOWN_INTERVAL);
      });
    }, 1000)
  };

  const endCountDown = () => {
    sendMessage(JSON.stringify({
      action: GameRoomState.END_COUNTER,
      roomID: gameRoom.roomID,
      numberOfQuestions,
    }));
  };

  const resetCountDown = () => {
    setTimeLeft(COUNT_DOWN_TIME);
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
    if (gameRoom.state === GameRoomState.COUNTING_DOWN && !counterOn) {
      startCountDown();
    }
  },[gameRoom.state, counterOn]);
  
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
    setGameRoom({...gameRoom, currentQuestionIndex: -1});
  }, []);

  useEffect(() => {
    if (lastMessage !== null) {
      const data = JSON.parse(lastMessage.data);
      setGameRoom({...gameRoom, ...data});
    }
  }, [lastMessage]);

  useEffect(() => {
    if (gameRoom.state === GameRoomState.END) {
      alert(`Game over! The winner is : ${gameRoom.winnerID}`);
      resetCountDown();
      setGameRoom({...gameRoom, state: GameRoomState.WAITING});
    }  
  });

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
