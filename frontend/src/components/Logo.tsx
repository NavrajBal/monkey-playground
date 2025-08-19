import React from 'react';
import { useTheme } from '../contexts/ThemeContext';
import './Logo.css';

interface LogoProps {
  size?: number;
  className?: string;
}

const Logo: React.FC<LogoProps> = ({ size = 32, className = '' }) => {
  const { theme } = useTheme();

  return (
    <div className={`logo-container ${className}`}>
      <img
        src="/monkey-logo.png"
        alt="Monkey Language Logo"
        className={`logo ${theme === 'dark' ? 'logo-dark' : 'logo-light'}`}
        style={{ width: size, height: size }}
      />
    </div>
  );
};

export default Logo;
