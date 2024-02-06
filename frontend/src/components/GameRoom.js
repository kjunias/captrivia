import { GameRoomState } from "../constants";
import useGameRoom from "../hooks/useGameRoom";
import LeaderBoard from "./LeaderBoard";
import Question from "./Question";

const GameRoom = () => {
  const {playerID, gameRoom, setGameRoom, startGame, timeLeft, resetTimer, setNumberOfQuestions, submitAnswer} = useGameRoom()

  const handleChange = (event) => {
    event.preventDefault();
    setNumberOfQuestions(parseInt(event.target.value));
  };

  const handleStartGame = (event) => {
    event.preventDefault();
    startGame();
  };


  const roomURL = `${window.location.origin}/join?roomID=${gameRoom?.roomID}`;

  if(gameRoom?.state === GameRoomState.PLAYING && gameRoom.questions?.length > 0 && gameRoom.currentQuestionIndex >= 0) {
    return (
    <>
        <Question
          questions={gameRoom.questions}
          currentQuestionIndex={gameRoom.currentQuestionIndex}
          score={gameRoom.scores[playerID]}
          submitAnswer={submitAnswer}
        />
        <LeaderBoard gameRoom={gameRoom}/>
      </>
    )
  }

  return (
    <>
      <h3>Game Room: {gameRoom?.roomID}</h3>
      {gameRoom?.state === GameRoomState.COUNTING_DOWN && (<h3>Countdown: {timeLeft/1000}</h3>)}
      <h4>Room URL: <a href={`${roomURL}`} target="_blank">{roomURL}</a></h4>
      <h4>PlayerID: {gameRoom?.playerID}</h4>
      {gameRoom?.playerID === gameRoom?.adminID && (
      <form>
        <label>
          Number of Questions:
          <input type="number" name="questionsnumber" onChange={handleChange} />
        </label> 
        <button onClick={handleStartGame}>Start Game</button>
      </form>)}
      <LeaderBoard gameRoom={gameRoom}/>
    </>
  );
};
export default GameRoom;
