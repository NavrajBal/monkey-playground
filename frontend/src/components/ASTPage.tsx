import React, { useState } from "react";
import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import MonacoEditor from "./MonacoEditor";
import ASTViewer from "./ASTViewer";
import SampleDropdown from "./SampleDropdown";
import { monkeyService, type ParseResponse } from "../services/monkeyService";
import { useTheme } from "../contexts/ThemeContext";
import { useCode } from "../contexts/CodeContext";
import { type CodeSample } from "../data/samples";
import "./ASTPage.css";

const ASTPage: React.FC = () => {
  const { code, setCode } = useCode();
  const [astData, setAstData] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string>("");
  const { theme } = useTheme();

  const handleCodeChange = (value: string | undefined) => {
    setCode(value || "");
  };

  const parseCode = async () => {
    if (!code.trim()) return;

    setIsLoading(true);
    setError("");

    try {
      const result: ParseResponse = await monkeyService.parse(code);

      if (result.error) {
        setError(result.error);
        setAstData(null);
      } else {
        setAstData(result.ast);
        setError("");
      }
    } catch (err) {
      setError(`Parse error: ${err}`);
      setAstData(null);
    } finally {
      setIsLoading(false);
    }
  };

  const clearAST = () => {
    setAstData(null);
    setError("");
  };

  const handleSelectSample = (sample: CodeSample) => {
    setCode(sample.code);
    setAstData(null);
    setError("");
  };

  return (
    <div className="ast-page">
      {/* Header */}
      <div className="header">
        <div className="header-left">
          <h1>AST Viewer</h1>
          <SampleDropdown onSelectSample={handleSelectSample} />
        </div>
        <div className="controls">
          <button
            onClick={parseCode}
            disabled={isLoading}
            className="parse-button"
          >
            {isLoading ? "Parsing..." : "Parse AST"}
          </button>
          <button onClick={clearAST} className="clear-button">
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

          {/* AST Visualization Panel */}
          <Panel defaultSize={60} minSize={40}>
            <div className="panel">
              <div className="panel-header">
                Abstract Syntax Tree
                {error && (
                  <span className="error-indicator">⚠️ Parse Error</span>
                )}
              </div>
              <div className="panel-content">
                {error ? (
                  <div className="error-content">
                    <div className="error-message">
                      <strong>Parse Error:</strong> {error}
                    </div>
                  </div>
                ) : (
                  <ASTViewer astData={astData} />
                )}
              </div>
            </div>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};

export default ASTPage;
