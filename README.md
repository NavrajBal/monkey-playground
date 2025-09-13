# 🐒 Monkey Playground

An interactive web-based playground for the **Monkey programming language** - a language designed for learning interpreters and compilers. Write, execute, and explore Monkey code right in your browser with real-time tokenization, AST visualization, and bytecode compilation. Find more information about MonkeyLang at https://monkeylang.org/.

For more details on the backend checkout my repo on monkey-lang, https://github.com/NavrajBal/monkey-lang.
## 📋 Table of Contents

- [🌐 Live Demo](#-live-demo)
- [🛠️ Technology Stack](#️-technology-stack)
- [🧠 Parser & Compiler](#-parser--compiler)
- [🎯 Features & Pages](#-features--pages)
- [🏗️ Architecture](#️-architecture)
- [🚀 Getting Started](#-getting-started)
- [🔧 Configuration](#-configuration)
- [📚 Monkey Language Features](#-monkey-language-features)
- [🚀 Deployment](#-deployment)
- [🚧 Work in Progress & Known Issues](#-work-in-progress--known-issues)
- [Acknowledgments](#acknowledgments)

## 🌐 Live Demo

**[Try it out @ ](https://monkey-playground.vercel.app)** *https://monkey-playground.vercel.app/*

### Video WalkThrough

## 🛠️ Technology Stack

### Frontend

- **React 19** with TypeScript
- **Vite** for build tooling and development
- **Monaco Editor** for code editing with syntax highlighting
- **React Flow** for AST visualization
- **React Router** for navigation
- **Tailwind CSS** equivalent styling
- **WebAssembly (WASM)** support for in-browser execution

### Backend

- **Go** HTTP server to be used locally
- **Vercel Functions** for serverless API endpoints
- RESTful API design

## 🧠 Parser & Compiler

This playground is powered by the **[monkey-lang](https://github.com/NavrajBal/monkey-lang)** implementation - a complete interpreter and compiler for the Monkey programming language. Which I made and extended features to taking inspiration from Thorsten Ball books "Writing an Interpreter in Go" and "Writing a Compiler in Go".

The implementation includes:

- **Lexical Analysis** - Tokenization of source code
- **Recursive Descent Parser** - AST generation
- **Tree-Walking Interpreter** - Direct AST evaluation
- **Bytecode Compiler** - Compilation to virtual machine instructions
- **Virtual Machine** - Bytecode execution engine

For more details checkout my repo on monkey-lang, https://github.com/NavrajBal/monkey-lang.
## 🎯 Features & Pages

### 🏠 Playground (Main REPL)

- Interactive Read-Eval-Print Loop
- Real-time code execution
- Sample code library with examples:
  - Fibonacci sequences
  - Factorial calculations
  - Array operations
  - Closures and higher-order functions
  - Hash maps
  - Conditional logic
- Dual backend support (API/WASM)
- Monaco editor with syntax highlighting

### 🔤 Tokenizer

- Real-time tokenization of Monkey source code
- Visual breakdown of code into tokens
- Token type identification (keywords, identifiers, operators, etc.)
- Helpful for understanding lexical analysis

### 🌳 AST Viewer

- Interactive Abstract Syntax Tree visualization
- Node-based graph representation using React Flow
- Expandable/collapsible tree structure
- Visual understanding of how code is parsed
- Color-coded node types for different AST elements

### ⚙️ Compiler (Coming Soon)

- Bytecode compilation visualization
- Instruction breakdown
- Constants pool inspection
- Virtual machine execution steps

## 🏗️ Architecture

### Frontend Structure

```
frontend/
├── src/
│   ├── components/          # React components
│   │   ├── Playground.tsx   # Main REPL interface
│   │   ├── TokenizerPage.tsx # Tokenization viewer
│   │   ├── ASTPage.tsx      # AST visualization
│   │   └── CompilerPage.tsx # Bytecode viewer
│   ├── contexts/           # React contexts
│   │   ├── CodeContext.tsx # Shared code state
│   │   ├── ThemeContext.tsx # Dark/light theme
│   │   └── ToastContext.tsx # Notifications
│   ├── services/           # API and WASM services
│   │   ├── apiService.ts   # HTTP API client
│   │   ├── wasmService.ts  # WebAssembly interface
│   │   └── monkeyService.ts # Unified service layer
│   └── wasm/              # Go-to-WASM compilation
│       └── main.go        # WASM entry point
├── api/                   # Vercel Functions
│   ├── tokenize.go       # Tokenization endpoint
│   ├── parse.go          # AST parsing endpoint
│   ├── compile.go        # Bytecode compilation
│   └── execute.go        # Code execution
└── public/
    └── monkey.wasm       # Compiled WebAssembly binary
```

### Backend Structure (only used locally)

```
backend/
├── main.go              # HTTP server entry point
├── api/
│   └── handlers.go      # API route handlers
└── go.mod              # Dependencies (monkey-lang)
```

## 🚀 Getting Started

### Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.22+ (for backend development)

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

### Backend Development

```bash
cd backend
go mod tidy
go run main.go
```

### Building WebAssembly

```bash
cd frontend
npm run build-wasm
```

## 🔧 Configuration

The playground supports two execution backends:

1. **API Backend** (Default) - Full-featured Go HTTP server
2. **WASM Backend** (Limited Features) - In-browser WebAssembly execution

Switch between backends using the toggle in the navigation bar or modify `frontend/src/config/config.ts`.


## 🚧 Work in Progress & Known Issues

### Work in Progress

- **Compiler Page**: Bytecode compilation visualization is currently under development
- **Enhanced Error Handling**: Improving error messages and debugging information
- **Performance Optimization**: Optimizing WASM execution for larger programs

### Known Issues

- **WASM Errors**: Code with Arrays, Hashmaps failing
- **WASM Tokenizer & AST Support**: WASM errors when trying to tokenize or display AST
- **Error Recovery**: Parser error recovery could be more robust


## Acknowledgments

**Thorsten Ball** for writing the excellent books on interpreters and compilers
