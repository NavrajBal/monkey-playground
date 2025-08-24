import React, { useState } from "react";
import "./TokenizerViewer.css";

interface TokenInfo {
  type: string;
  literal: string;
  position: number;
}

interface TokenizerViewerProps {
  tokens: TokenInfo[];
  code: string;
}

const TokenizerViewer: React.FC<TokenizerViewerProps> = ({ tokens, code }) => {
  const [selectedTokenIndex, setSelectedTokenIndex] = useState<number | null>(
    null
  );
  const [hoveredTokenIndex, setHoveredTokenIndex] = useState<number | null>(
    null
  );

  const getTokenColor = (tokenType: string): string => {
    switch (tokenType) {
      case "LET":
      case "IF":
      case "ELSE":
      case "RETURN":
      case "FN":
      case "TRUE":
      case "FALSE":
        return "#8b5cf6"; // Purple - Keywords
      case "IDENT":
        return "#06b6d4"; // Cyan - Identifiers
      case "INT":
        return "#10b981"; // Green - Numbers
      case "STRING":
        return "#f59e0b"; // Yellow - Strings
      case "+":
      case "-":
      case "*":
      case "/":
      case "=":
      case "==":
      case "!=":
      case "<":
      case ">":
      case "!":
        return "#ef4444"; // Red - Operators
      case "(":
      case ")":
      case "{":
      case "}":
      case "[":
      case "]":
        return "#ec4899"; // Pink - Delimiters
      case ";":
      case ",":
        return "#6b7280"; // Gray - Punctuation
      case "EOF":
        return "#64748b"; // Slate - Special
      case "ILLEGAL":
        return "#dc2626"; // Dark red - Errors
      default:
        return "#94a3b8"; // Light gray - Unknown
    }
  };

  const getTokenDisplayName = (tokenType: string): string => {
    switch (tokenType) {
      case "IDENT":
        return "Identifier";
      case "INT":
        return "Integer";
      case "STRING":
        return "String";
      case "LET":
        return "Let Keyword";
      case "FN":
        return "Function Keyword";
      case "IF":
        return "If Keyword";
      case "ELSE":
        return "Else Keyword";
      case "RETURN":
        return "Return Keyword";
      case "TRUE":
        return "Boolean True";
      case "FALSE":
        return "Boolean False";
      case "EOF":
        return "End of File";
      case "ILLEGAL":
        return "Illegal Token";
      default:
        return tokenType;
    }
  };

  const renderHighlightedCode = () => {
    if (!tokens.length || !code) {
      return (
        <div className="highlighted-code">{code || "No code to tokenize"}</div>
      );
    }

    const elements: React.ReactElement[] = [];
    let lastPosition = 0;

    tokens.forEach((token, index) => {
      // Add any text between tokens
      if (token.position > lastPosition) {
        const betweenText = code.slice(lastPosition, token.position);
        if (betweenText) {
          elements.push(
            <span key={`between-${index}`} className="token-whitespace">
              {betweenText}
            </span>
          );
        }
      }

      // Add the token
      const isSelected = selectedTokenIndex === index;
      const isHovered = hoveredTokenIndex === index;
      const tokenEnd = token.position + token.literal.length;

      elements.push(
        <span
          key={`token-${index}`}
          className={`token ${isSelected ? "selected" : ""} ${
            isHovered ? "hovered" : ""
          }`}
          style={{
            backgroundColor:
              getTokenColor(token.type) +
              (isSelected ? "ff" : isHovered ? "cc" : "99"),
            color: "white",
          }}
          onClick={() => setSelectedTokenIndex(index)}
          onMouseEnter={() => setHoveredTokenIndex(index)}
          onMouseLeave={() => setHoveredTokenIndex(null)}
        >
          {token.literal}
        </span>
      );

      lastPosition = tokenEnd;
    });

    // Add any remaining text
    if (lastPosition < code.length) {
      const remainingText = code.slice(lastPosition);
      if (remainingText) {
        elements.push(
          <span key="remaining" className="token-whitespace">
            {remainingText}
          </span>
        );
      }
    }

    return <div className="highlighted-code">{elements}</div>;
  };

  if (!tokens.length) {
    return (
      <div className="tokenizer-viewer-empty">
        <p>No tokens to display. Tokenize some code first!</p>
      </div>
    );
  }

  return (
    <div className="tokenizer-viewer">
      {/* Highlighted Code Section */}
      <div className="highlighted-section">
        <div className="section-header">
          <h3>Tokenized Code</h3>
          <span className="token-count">{tokens.length} tokens</span>
        </div>
        <div className="highlighted-container">{renderHighlightedCode()}</div>
      </div>

      {/* Token List Section */}
      <div className="token-list-section">
        <div className="section-header">
          <h3>Token Details</h3>
          {selectedTokenIndex !== null && (
            <button
              className="clear-selection"
              onClick={() => setSelectedTokenIndex(null)}
            >
              Clear Selection
            </button>
          )}
        </div>
        <div className="token-list">
          {tokens.map((token, index) => (
            <div
              key={index}
              className={`token-item ${
                selectedTokenIndex === index ? "selected" : ""
              } ${hoveredTokenIndex === index ? "hovered" : ""}`}
              onClick={() => setSelectedTokenIndex(index)}
              onMouseEnter={() => setHoveredTokenIndex(index)}
              onMouseLeave={() => setHoveredTokenIndex(null)}
            >
              <div className="token-item-header">
                <span
                  className="token-type-badge"
                  style={{ backgroundColor: getTokenColor(token.type) }}
                >
                  {getTokenDisplayName(token.type)}
                </span>
                <span className="token-position">@{token.position}</span>
              </div>
              <div className="token-literal">"{token.literal}"</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TokenizerViewer;
