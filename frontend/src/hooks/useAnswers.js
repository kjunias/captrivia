import { useContext } from "react";
import { API_BASE } from "../config";
import { GameContext } from "../contexts/GameContext";
import { QuestionContext } from "../contexts/QuestionContext";

const useAnswers = () => {
  const { setError, setLoading, gameSession, setGameSession } =
    useContext(GameContext);
  const {
    questions,
    setQuestions,
    currentQuestionIndex,
    setCurrentQuestionIndex,
    setScore,
  } = useContext(QuestionContext);

  const submitAnswer = async (index) => {
    // We are submitting the index
    setLoading(true);
    const currentQuestion = questions[currentQuestionIndex];
    try {
      const res = await fetch(`${API_BASE}/answer`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sessionId: gameSession,
          questionId: currentQuestion.id, // field name is "id", not "questionId"
          answer: index,
        }),
      });
      const data = await res.json();
      if (data.correct) {
        setScore(data.currentScore); // Update score from server's response
      }
      if (currentQuestionIndex < questions.length - 1) {
        setCurrentQuestionIndex(currentQuestionIndex + 1);
      } else {
        endGame();
      }
    } catch (err) {
      setError("Failed to submit answer.");
    }
    setLoading(false);
  };

  const endGame = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/game/end`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sessionId: gameSession, // need to provide the sessionId
        }),
      });
      const data = await res.json();
      alert(`Game over! Your score: ${data.finalScore}`); // Use the finalScore from the response
      setGameSession(null);
      setQuestions([]);
      setCurrentQuestionIndex(0);
      setScore(0);
    } catch (err) {
      setError("Failed to end game.");
    }
    setLoading(false);
  };

  return {
    submitAnswer,
  };
};
export default useAnswers;
