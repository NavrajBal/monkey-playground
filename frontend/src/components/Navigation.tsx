import React from "react";
import { Link, useLocation } from "react-router-dom";
import { Sun, Moon, Play, Code, TreePine, Cpu } from "lucide-react";
import { useTheme } from "../contexts/ThemeContext";
import Logo from "./Logo";
import { BackendToggle } from "./BackendToggle";
import "./Navigation.css";

const Navigation: React.FC = () => {
  const { theme, toggleTheme } = useTheme();
  const location = useLocation();

  const navItems = [
    { path: "/", label: "REPL", icon: Play },
    { path: "/tokenizer", label: "Tokenizer", icon: Code },
    { path: "/ast", label: "AST Viewer", icon: TreePine },
    { path: "/compiler", label: "Compiler", icon: Cpu },
  ];

  return (
    <nav className="navigation">
      <div className="nav-brand">
        <Logo size={28} className="nav-logo" />
        <span className="nav-title">Monkey Playground</span>
      </div>

      <div className="nav-items">
        {navItems.map(({ path, label, icon: Icon }) => (
          <Link
            key={path}
            to={path}
            className={`nav-item ${location.pathname === path ? "active" : ""}`}
          >
            <Icon size={18} />
            <span>{label}</span>
          </Link>
        ))}
      </div>

      <div className="nav-controls">
        <BackendToggle />
        <button
          onClick={toggleTheme}
          className="theme-toggle"
          aria-label="Toggle theme"
        >
          {theme === "dark" ? <Sun size={20} /> : <Moon size={20} />}
        </button>
      </div>
    </nav>
  );
};

export default Navigation;
