import React from 'react';
import { HelpCircle } from 'lucide-react';
import './HelpButton.css';

interface HelpButtonProps {
  onClick: () => void;
}

export const HelpButton: React.FC<HelpButtonProps> = ({ onClick }) => {
  return (
    <button
      onClick={onClick}
      className="help-button"
      aria-label="Show help and information"
      title="Show help and information"
    >
      <HelpCircle size={20} />
    </button>
  );
};
