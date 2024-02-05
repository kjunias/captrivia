
const LeaderBoard = ({gameRoom}) => {
  return (
    <table>
    <thead>
      <tr><th>Leader Board</th></tr>
    </thead>
    <tbody>
      {
      Object.entries(gameRoom?.scores || {})
       .sort(([,a], [,b]) => b-a)
       .map((pId, i) => {
        return (
          <tr key={pId}>
            <td>{i + 1}: {pId[0]}</td>
            <td/>
            <td>&nbsp;{pId[1]}</td>
          </tr>
        )
      })
      }
    </tbody>
  </table>
);
};

export default LeaderBoard;