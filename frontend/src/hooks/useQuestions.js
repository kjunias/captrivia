import { useContext, useEffect } from "react";
import { GameContext } from "../contexts/GameContext";
import { API_BASE } from "../config";
import { QuestionContext } from "../contexts/QuestionContext";

const useQuestions = () => {
  const { loading, setLoading, gameSession, setError } =
    useContext(GameContext);
  const { questions, setQuestions, currentQuestionIndex, score, setScore } =
    useContext(QuestionContext);

  useEffect(() => {
    if (gameSession && !loading && questions.length <= 0) {
      const fetchQuestions = async function () {
        setLoading(true);
        try {
          const res = await fetch(`${API_BASE}/questions`);
          const data = await res.json();
          setQuestions(data);
        } catch (err) {
          setError("Failed to fetch questions.");
        }
        setLoading(false);
      };
      fetchQuestions();
    }
  }, []);

  return {
    questions,
    setQuestions,
    currentQuestionIndex,
    score,
    setScore,
  };
};

export default useQuestions;
