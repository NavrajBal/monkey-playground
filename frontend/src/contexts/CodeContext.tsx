import React, {
  createContext,
  useContext,
  useState,
  type ReactNode,
} from "react";

interface CodeContextType {
  code: string;
  setCode: (code: string) => void;
}

const CodeContext = createContext<CodeContextType | undefined>(undefined);

interface CodeProviderProps {
  children: ReactNode;
}

const defaultCode = `let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      return 1;
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};

fibonacci(10);`;

export const CodeProvider: React.FC<CodeProviderProps> = ({ children }) => {
  const [code, setCode] = useState(defaultCode);

  return (
    <CodeContext.Provider value={{ code, setCode }}>
      {children}
    </CodeContext.Provider>
  );
};

export const useCode = (): CodeContextType => {
  const context = useContext(CodeContext);
  if (context === undefined) {
    throw new Error("useCode must be used within a CodeProvider");
  }
  return context;
};
