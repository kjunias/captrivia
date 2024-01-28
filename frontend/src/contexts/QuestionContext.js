import { createContext, useState } from "react";

export const QuestionContext = createContext({});

export const QuestionProvider = ({ children }) => {
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [score, setScore] = useState(0);

  return (
    <QuestionContext.Provider
      value={{
        questions,
        setQuestions,
        currentQuestionIndex,
        setCurrentQuestionIndex,
        score,
        setScore,
      }}
    >
      {children}
    </QuestionContext.Provider>
  );
};
