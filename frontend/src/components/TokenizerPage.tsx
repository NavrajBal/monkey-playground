import React, { useState } from "react";
import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import MonacoEditor from "./MonacoEditor";
import TokenizerViewer from "./TokenizerViewer";
import SampleDropdown from "./SampleDropdown";
import {
  monkeyService,
  type TokenizeResponse,
  type TokenInfo,
} from "../services/monkeyService";
import { useTheme } from "../contexts/ThemeContext";
import { useCode } from "../contexts/CodeContext";
import { type CodeSample } from "../data/samples";
import "./TokenizerPage.css";

const TokenizerPage: React.FC = () => {
  const { code, setCode } = useCode();
  const [tokens, setTokens] = useState<TokenInfo[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string>("");
  const { theme } = useTheme();

  const handleCodeChange = (value: string | undefined) => {
    setCode(value || "");
  };

  const tokenizeCode = async () => {
    if (!code.trim()) return;

    setIsLoading(true);
    setError("");

    try {
      const result: TokenizeResponse = await monkeyService.tokenize(code);

      if (result.error) {
        setError(result.error);
        setTokens([]);
      } else {
        setTokens(result.tokens || []);
        setError("");
      }
    } catch (err) {
      setError(`Tokenization error: ${err}`);
      setTokens([]);
    } finally {
      setIsLoading(false);
    }
  };

  const clearTokens = () => {
    setTokens([]);
    setError("");
  };

  const handleSelectSample = (sample: CodeSample) => {
    setCode(sample.code);
    setTokens([]);
    setError("");
  };

  return (
    <div className="tokenizer-page">
      {/* Header */}
      <div className="header">
        <div className="header-left">
          <h1>Tokenizer</h1>
          <SampleDropdown onSelectSample={handleSelectSample} />
        </div>
        <div className="controls">
          <button
            onClick={tokenizeCode}
            disabled={isLoading}
            className="tokenize-button"
          >
            {isLoading ? "Tokenizing..." : "Tokenize"}
          </button>
          <button onClick={clearTokens} className="clear-button">
            Clear
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="main-content">
        <PanelGroup direction="horizontal">
          {/* Code Editor Panel */}
          <Panel defaultSize={40} minSize={30}>
            <div className="panel">
              <div className="panel-header">Code Editor</div>
              <div className="panel-content">
                <MonacoEditor
                  value={code}
                  onChange={handleCodeChange}
                  height="100%"
                  theme={theme === "dark" ? "vs-dark" : "light"}
                />
              </div>
            </div>
          </Panel>

          <PanelResizeHandle className="resize-handle" />

          {/* Tokenizer Visualization Panel */}
          <Panel defaultSize={60} minSize={40}>
            <div className="panel">
              <div className="panel-header">
                Token Visualization
                {error && (
                  <span className="error-indicator">⚠️ Tokenization Error</span>
                )}
              </div>
              <div className="panel-content">
                {error ? (
                  <div className="error-content">
                    <div className="error-message">
                      <strong>Tokenization Error:</strong> {error}
                    </div>
                  </div>
                ) : (
                  <TokenizerViewer tokens={tokens} code={code} />
                )}
              </div>
            </div>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};

export default TokenizerPage;
