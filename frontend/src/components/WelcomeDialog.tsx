import React, { useEffect, useState } from 'react';
import { X, AlertTriangle, Zap, Code, Server } from 'lucide-react';
import './WelcomeDialog.css';

interface WelcomeDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export const WelcomeDialog: React.FC<WelcomeDialogProps> = ({ isOpen, onClose }) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setIsVisible(true);
    }
  }, [isOpen]);

  const handleClose = () => {
    setIsVisible(false);
    setTimeout(onClose, 300); // Wait for fade out animation
  };

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      handleClose();
    }
  };

  if (!isOpen) return null;

  return (
    <div 
      className={`welcome-overlay ${isVisible ? 'welcome-overlay-visible' : ''}`}
      onClick={handleOverlayClick}
    >
      <div className={`welcome-dialog ${isVisible ? 'welcome-dialog-visible' : ''}`}>
        <div className="welcome-header">
          <div className="welcome-title">
            <Code size={24} />
            Welcome to Monkey Playground!
          </div>
          <button
            onClick={handleClose}
            className="welcome-close"
            aria-label="Close welcome dialog"
          >
            <X size={20} />
          </button>
        </div>

        <div className="welcome-content">
          <div className="welcome-section">
            <h3>🎉 What is this?</h3>
            <p>
              An interactive playground for the <strong>Monkey programming language</strong> - 
              a language designed for learning interpreters and compilers. You can write, 
              execute, and explore Monkey code right in your browser!
            </p>
          </div>

          <div className="welcome-section">
            <h3>🚀 Features</h3>
            <ul>
              <li><strong>REPL</strong> - Interactive code execution</li>
              <li><strong>Tokenizer</strong> - See how code is broken into tokens</li>
              <li><strong>AST Viewer</strong> - Visualize the Abstract Syntax Tree</li>
              <li><strong>Compiler</strong> - View bytecode compilation (coming soon)</li>
            </ul>
          </div>

          <div className="welcome-section warning-section">
            <div className="warning-header">
              <AlertTriangle size={20} />
              <h3>⚠️ Important Notes</h3>
            </div>
            <div className="backend-info">
              <div className="backend-option">
                <Server size={16} />
                <div>
                  <strong>API Backend (Default)</strong>
                  <p>Uses Go server - fully featured and stable</p>
                </div>
              </div>
              <div className="backend-option">
                <Zap size={16} />
                <div>
                  <strong>WASM Backend</strong>
                  <p>Runs in browser - faster but some features still in development</p>
                </div>
              </div>
            </div>
            <p className="warning-text">
              You can switch between backends using the toggle in the navigation bar. 
              If you encounter issues with WASM, try switching to the API backend.
            </p>
          </div>

          <div className="welcome-section">
            <h3>📚 Quick Start</h3>
            <p>
              Try the sample code in the dropdown, or write your own Monkey code. 
              The language supports functions, arrays, hash maps, and more!
            </p>
          </div>
        </div>

        <div className="welcome-footer">
          <button onClick={handleClose} className="welcome-button">
            Let's Code! 🚀
          </button>
        </div>
      </div>
    </div>
  );
};
