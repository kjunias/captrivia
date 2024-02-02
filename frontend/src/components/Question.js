const Question = ({questions, currentQuestionIndex, submitAnswer, score}) => {
  if (!(questions && currentQuestionIndex)) {
    return null
  }
  return (
    <div>
      <h3>{questions[currentQuestionIndex]?.questionText}</h3>
      {questions[currentQuestionIndex]?.options.map((option, index) => (
        <button
          key={index} // Key should be unique for each child in a list, use index as the key
          onClick={() => submitAnswer(index)} // Pass index instead of option
          className="option-button"
        >
          {option}
        </button>
      ))}
      <p className="score">Score: {score}</p>
    </div>
  );
};

export default Question;
