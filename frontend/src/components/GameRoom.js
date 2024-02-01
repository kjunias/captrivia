import useGameRoom from "../hooks/useGameRoom";
import useWebSocket, { ReadyState } from "react-use-websocket"

const GameRoom = () => {
  const {gameRoom, startGame, timeLeft, setNumberOfQuestions} = useGameRoom()

  const handleChange = (event) => {
    event.preventDefault();
    setNumberOfQuestions(parseInt(event.target.value));
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    startGame();
  };

  const roomURL = `${window.location.origin}/join?roomID=${gameRoom?.roomID}`;

  return (
    <>
      <h3>Game Room: {gameRoom?.roomID}</h3>
      {gameRoom?.isCountingDown && (<h3>Countdown: {timeLeft}</h3>)}
      <h4>Room URL: <a href={`${roomURL}`} target="_blank">{roomURL}</a></h4>
      <h4>PlayerID: {gameRoom?.playerID}</h4>
      {gameRoom?.playerID === gameRoom?.adminID && (
      <form>
        <label>
          Number of Questions:
          <input type="number" name="questionsnumber" onChange={handleChange} />
        </label> 
        <button onClick={handleSubmit}>Start Game</button>
      </form>)}
      <table>
        <thead>
          <tr><th>Players</th></tr>
        </thead>
        <tbody>
          {Object.keys(gameRoom?.scores || {}).map((pId, i) => {
            return (
              <tr key={pId}>
                <td>{i + 1}: {pId}</td>
                <td>{gameRoom.scores[pId]}</td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </>
  );
};
export default GameRoom;
