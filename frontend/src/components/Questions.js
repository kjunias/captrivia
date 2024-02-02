import { useContext } from "react";
import { GameContext } from "../contexts/GameContext";
import { QuestionsProvider } from "../contexts/QuestionContext";
import useAnswers from "../hooks/useAnswers";
import useQuestions from "../hooks/useQuestions";
import Question from "./Question";
// import Question from "./Question";

const Questions = () => {
  const { loading, error } = useContext(GameContext);
  const { questions, currentQuestionIndex, score } = useQuestions();
  const { submitAnswer } = useAnswers();

  if (error) return <div className="error">Error: {error}</div>;
  if (loading) return <div className="loading">Loading...</div>;

  debugger;
  return (
    <Question
      questions={questions}
      currentQuestionIndex={currentQuestionIndex}
      score={score}
      submitAnswer={submitAnswer}
    />
  );
};

const QuestionsContainer = () => {
  return (
    <QuestionsProvider>
      <Questions />
    </QuestionsProvider>
  );
};
export default QuestionsContainer;
