import React, { useState } from "react";
import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import MonacoEditor from "./MonacoEditor";
import SampleDropdown from "./SampleDropdown";
import { monkeyService, type ExecuteResponse } from "../services/monkeyService";
import { useTheme } from "../contexts/ThemeContext";
import { useCode } from "../contexts/CodeContext";
import { type CodeSample } from "../data/samples";
import "./Playground.css";

const Playground: React.FC = () => {
  const { code, setCode } = useCode();
  const [output, setOutput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { theme } = useTheme();

  const handleCodeChange = (value: string | undefined) => {
    setCode(value || "");
  };

  const executeCode = async () => {
    if (!code.trim()) return;

    setIsLoading(true);
    setOutput("Executing...");

    try {
      const result: ExecuteResponse = await monkeyService.execute(code);

      if (!result) {
        setOutput("Error: No response from WASM");
        return;
      }

      if (result.error) {
        setOutput(`Error: ${result.error}`);
      } else {
        let output = "";
        if (result.output) {
          output += result.output;
        }
        if (result.result && result.result !== "null") {
          if (output) output += "\n";
          output += `=> ${result.result}`;
        }
        setOutput(output || result.result || "");
      }
    } catch (error) {
      setOutput(`Error: ${error}`);
    } finally {
      setIsLoading(false);
    }
  };

  const clearOutput = () => {
    setOutput("");
  };

  const handleSelectSample = (sample: CodeSample) => {
    setCode(sample.code);
    setOutput("");
  };

  return (
    <div className="playground">
      {/* Header */}
      <div className="header">
        <div className="header-left">
          <h1>REPL</h1>
          <SampleDropdown onSelectSample={handleSelectSample} />
        </div>
        <div className="controls">
          <button
            onClick={executeCode}
            disabled={isLoading}
            className="run-button"
          >
            {isLoading ? "Running..." : "Run Code"}
          </button>
          <button onClick={clearOutput} className="clear-button">
            Clear
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="main-content">
        <PanelGroup direction="horizontal">
          {/* Code Editor Panel */}
          <Panel defaultSize={60} minSize={30}>
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

          {/* Output Panel */}
          <Panel defaultSize={40} minSize={20}>
            <div className="panel">
              <div className="panel-header">Output</div>
              <div className="output-content">
                {output ? (
                  <pre className="output-text">{output}</pre>
                ) : (
                  <div className="output-placeholder">
                    Run your Monkey code to see the output here...
                  </div>
                )}
              </div>
            </div>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};

export default Playground;
