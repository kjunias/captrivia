import { useContext } from "react";
import { GameContext } from "../contexts/GameContext";
import useGameRoom from "../hooks/useGameRoom";

const GameRoom = () => {
  const { gameRoom, playerID} =
    useContext(GameContext);
  const {update} = useGameRoom()

  return (
    <>
      <h3>Game Room: {gameRoom}</h3>
      <h4>PlayerID: {playerID}</h4>
      <table>
        <thead>
          <tr><th>Players</th></tr>
        </thead>
        <tbody>
          {update?.scores?.map((p, x) => {
            console.log("p:", p)
            console.log("x:", x)
            return (
              <tr>
                <td>{p}</td>
                <td>0</td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </>
  )
};
export default GameRoom;
