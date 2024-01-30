import { useContext, useEffect, useState } from "react";
import { GameContext } from "../contexts/GameContext";
import { API_BASE } from "../config";

const useGameRoom = () => {
  const { gameRoom, playerID } = useContext(GameContext);
  const [update, setUpdate] = useState(null);

  useEffect(() => {
      const fetchUpdates = async function () {
        try {
          const res = await fetch(`${API_BASE}/gameroom/update?roomID=${gameRoom}&playerID=${playerID}`,{
            method: "GET",
            headers: {
              "Content-Type": "application/json",
              "Accept": "application/json",
            },
          });
          const data =  await new Response(res.body).text()
          console.log("Update res:", data);
          //const data = await res.json();
          console.log("Update status:", res.status);
          console.log("Update data:", data);
          setUpdate(data);
          fetchUpdates();
        } catch (err) {
          console.error("Failed to fetch updates: ", err);
        }
      };
      fetchUpdates();
  }, []);

  return {
    update
  };
};

export default useGameRoom;
