import { useContext } from "react";
import { GameContext } from "../contexts/GameContext";
import { QuestionProvider } from "../contexts/QuestionContext";
import useAnswers from "../hooks/useAnswers";
import useQuestions from "../hooks/useQuestions";

const Question = () => {
  const { loading, error } = useContext(GameContext);
  const { submitAnswer } = useAnswers();
  const { questions, currentQuestionIndex, score } = useQuestions();

  if (error) return <div className="error">Error: {error}</div>;
  if (loading) return <div className="loading">Loading...</div>;

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

const QuestionContainer = () => {
  return (
    <QuestionProvider>
      <Question />
    </QuestionProvider>
  );
};
export default QuestionContainer;
