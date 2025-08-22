// Configuration for the Monkey Playground
export const config = {
  // Toggle between 'api' and 'wasm' backends
  // 'api' - Uses Go HTTP server backend
  // 'wasm' - Uses WebAssembly in browser
  backend: "api" as "api" | "wasm",

  // API configuration (when using 'api' backend)
  apiUrl: import.meta.env.PROD
    ? "/api" // Vercel Functions will be available at /api/*
    : "http://localhost:8080",

  // Development settings
  enableDebugLogs: true,
};

// Helper function to check current backend
export const isUsingWasm = () => config.backend === "wasm";
export const isUsingApi = () => config.backend === "api";
