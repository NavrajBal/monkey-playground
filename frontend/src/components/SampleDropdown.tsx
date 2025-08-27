import React, { useState } from "react";
import { ChevronDown, BookOpen } from "lucide-react";
import { codeSamples, type CodeSample } from "../data/samples";
import "./SampleDropdown.css";

interface SampleDropdownProps {
  onSelectSample: (sample: CodeSample) => void;
}

const SampleDropdown: React.FC<SampleDropdownProps> = ({ onSelectSample }) => {
  const [isOpen, setIsOpen] = useState(false);

  const handleSelectSample = (sample: CodeSample) => {
    onSelectSample(sample);
    setIsOpen(false);
  };

  return (
    <div className="sample-dropdown">
      <button onClick={() => setIsOpen(!isOpen)} className="sample-trigger">
        <BookOpen size={18} />
        <span>Examples</span>
        <ChevronDown size={16} className={`chevron ${isOpen ? "open" : ""}`} />
      </button>

      {isOpen && (
        <div className="sample-menu">
          <div className="sample-header">
            <span>Sample Code</span>
          </div>
          {codeSamples.map((sample) => (
            <button
              key={sample.id}
              onClick={() => handleSelectSample(sample)}
              className="sample-item"
            >
              <div className="sample-title">{sample.title}</div>
              <div className="sample-description">{sample.description}</div>
            </button>
          ))}
        </div>
      )}

      {isOpen && (
        <div className="sample-overlay" onClick={() => setIsOpen(false)} />
      )}
    </div>
  );
};

export default SampleDropdown;
