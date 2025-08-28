import React from 'react';
import { Settings, Server, Globe } from 'lucide-react';
import { config } from '../config/config';
import { useToast } from '../contexts/ToastContext';
import './BackendToggle.css';

interface BackendToggleProps {
  className?: string;
}

export const BackendToggle: React.FC<BackendToggleProps> = ({ className = '' }) => {
  const [backend, setBackend] = React.useState<'api' | 'wasm'>(config.backend);
  const { showToast } = useToast();

  const handleToggle = (newBackend: 'api' | 'wasm') => {
    config.backend = newBackend;
    setBackend(newBackend);
    
    // Show toast notification
    if (newBackend === 'wasm') {
      showToast(
        '⚡ Switched to WebAssembly - Fast in-browser execution!',
        'info',
        4000
      );
    } else {
      showToast(
        '🖥️ Switched to API Backend - Full-featured Go server',
        'success',
        4000
      );
    }
  };

  return (
    <div className={`backend-toggle ${className}`}>
      <div className="backend-toggle-label">
        <Settings size={16} />
        Backend:
      </div>
      <div className="backend-toggle-buttons">
        <button
          className={`backend-button ${backend === 'api' ? 'active' : ''}`}
          onClick={() => handleToggle('api')}
          title="Use Go HTTP server backend"
        >
          <Server size={14} />
          API
        </button>
        <button
          className={`backend-button ${backend === 'wasm' ? 'active' : ''}`}
          onClick={() => handleToggle('wasm')}
          title="Use WebAssembly in-browser execution"
        >
          <Globe size={14} />
          WASM
        </button>
      </div>
    </div>
  );
};
