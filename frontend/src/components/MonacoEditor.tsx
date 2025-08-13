import React from "react";
import Editor from "@monaco-editor/react";

interface MonacoEditorProps {
  value: string;
  onChange: (value: string | undefined) => void;
  language?: string;
  theme?: string;
  height?: string;
}

const MonacoEditor: React.FC<MonacoEditorProps> = ({
  value,
  onChange,
  language = "javascript", // We'll customize this for Monkey language later
  theme = "vs-dark",
  height = "400px",
}) => {
  return (
    <Editor
      height={height}
      language={language}
      theme={theme}
      value={value}
      onChange={onChange}
      options={{
        minimap: { enabled: false },
        fontSize: 14,
        lineNumbers: "on",
        roundedSelection: false,
        scrollBeyondLastLine: false,
        readOnly: false,
        automaticLayout: true,
        // Disable error squiggles/underlining
        "semanticHighlighting.enabled": false,
        renderValidationDecorations: "off",
      }}
    />
  );
};

export default MonacoEditor;
