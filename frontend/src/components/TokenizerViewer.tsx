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

    // Process tokens in order, finding them sequentially in the code
    // This ensures correct ordering and prevents misplacement
    const elements: React.ReactElement[] = [];
    let codePos = 0;

    tokens.forEach((token, tokenIdx) => {
      // Search for this token literal starting from where we've processed so far
      // This ensures we find tokens in the correct order
      const remainingCode = code.slice(codePos);
      const tokenIndexInCode = remainingCode.indexOf(token.literal);

      if (tokenIndexInCode === -1) {
        // Token not found - this shouldn't happen, but handle gracefully
        // Try to find it near the expected position as fallback
        const searchStart = Math.max(0, token.position - 20);
        const searchEnd = Math.min(
          code.length,
          token.position + token.literal.length + 20
        );
        const searchText = code.slice(searchStart, searchEnd);
        const fallbackIndex = searchText.indexOf(token.literal);

        if (fallbackIndex !== -1) {
          const actualTokenStart = searchStart + fallbackIndex;

          // Add whitespace before this token
          if (actualTokenStart > codePos) {
            const whitespace = code.slice(codePos, actualTokenStart);
            if (whitespace) {
              elements.push(
                <span key={`ws-${tokenIdx}`} className="token-whitespace">
                  {whitespace}
                </span>
              );
            }
          }

          // Add the token
          const actualTokenEnd = actualTokenStart + token.literal.length;
          const tokenText = code.slice(actualTokenStart, actualTokenEnd);

          const isSelected = selectedTokenIndex === tokenIdx;
          const isHovered = hoveredTokenIndex === tokenIdx;

          if (tokenText === token.literal) {
            elements.push(
              <span
                key={`token-${tokenIdx}`}
                className={`token ${isSelected ? "selected" : ""} ${
                  isHovered ? "hovered" : ""
                }`}
                style={{
                  backgroundColor:
                    getTokenColor(token.type) +
                    (isSelected ? "ff" : isHovered ? "cc" : "99"),
                  color: "white",
                }}
                onClick={() => setSelectedTokenIndex(tokenIdx)}
                onMouseEnter={() => setHoveredTokenIndex(tokenIdx)}
                onMouseLeave={() => setHoveredTokenIndex(null)}
              >
                {tokenText}
              </span>
            );
            codePos = actualTokenEnd;
          }
        }
        return; // Skip this token if we can't find it
      }

      // Found the token at the correct position
      const actualTokenStart = codePos + tokenIndexInCode;

      // Add whitespace before this token
      if (actualTokenStart > codePos) {
        const whitespace = code.slice(codePos, actualTokenStart);
        if (whitespace) {
          elements.push(
            <span key={`ws-${tokenIdx}`} className="token-whitespace">
              {whitespace}
            </span>
          );
        }
      }

      // Add the token
      const actualTokenEnd = actualTokenStart + token.literal.length;
      const tokenText = code.slice(actualTokenStart, actualTokenEnd);

      // Verify the token text matches (safety check)
      if (tokenText !== token.literal) {
        console.warn(
          `Token mismatch at position ${actualTokenStart}: expected "${token.literal}", got "${tokenText}"`
        );
      }

      const isSelected = selectedTokenIndex === tokenIdx;
      const isHovered = hoveredTokenIndex === tokenIdx;

      elements.push(
        <span
          key={`token-${tokenIdx}`}
          className={`token ${isSelected ? "selected" : ""} ${
            isHovered ? "hovered" : ""
          }`}
          style={{
            backgroundColor:
              getTokenColor(token.type) +
              (isSelected ? "ff" : isHovered ? "cc" : "99"),
            color: "white",
          }}
          onClick={() => setSelectedTokenIndex(tokenIdx)}
          onMouseEnter={() => setHoveredTokenIndex(tokenIdx)}
          onMouseLeave={() => setHoveredTokenIndex(null)}
        >
          {tokenText}
        </span>
      );

      codePos = actualTokenEnd;
    });

    // Don't render anything after the last token - this prevents duplicate text
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
