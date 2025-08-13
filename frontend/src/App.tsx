import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "./contexts/ThemeContext";
import { CodeProvider } from "./contexts/CodeContext";
import { ToastProvider } from "./contexts/ToastContext";
import Navigation from "./components/Navigation";
import Playground from "./components/Playground";
import TokenizerPage from "./components/TokenizerPage";
import ASTPage from "./components/ASTPage";
import CompilerPage from "./components/CompilerPage";
import { WelcomeDialog } from "./components/WelcomeDialog";
import { HelpButton } from "./components/HelpButton";
import "./App.css";

function App() {
  const [showWelcome, setShowWelcome] = useState(false);

  useEffect(() => {
    // Show welcome dialog on each reload
    setShowWelcome(true);
  }, []);

  const handleWelcomeClose = () => {
    setShowWelcome(false);
  };

  const handleShowWelcome = () => {
    setShowWelcome(true);
  };

  return (
    <ThemeProvider>
      <CodeProvider>
        <ToastProvider>
          <Router>
            <div className="app">
              <Navigation />
              <main className="main">
                <Routes>
                  <Route path="/" element={<Playground />} />
                  <Route path="/tokenizer" element={<TokenizerPage />} />
                  <Route path="/ast" element={<ASTPage />} />
                  <Route path="/compiler" element={<CompilerPage />} />
                </Routes>
              </main>
              <WelcomeDialog
                isOpen={showWelcome}
                onClose={handleWelcomeClose}
              />
              <HelpButton onClick={handleShowWelcome} />
            </div>
          </Router>
        </ToastProvider>
      </CodeProvider>
    </ThemeProvider>
  );
}

export default App;
