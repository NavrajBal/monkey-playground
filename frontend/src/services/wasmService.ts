// WASM Service - replaces API calls with direct WASM function calls

interface TokenInfo {
  type: string;
  literal: string;
  position: number;
}

interface ParsedAST {
  type: string;
  [key: string]: any;
}

interface TokenizeResponse {
  tokens?: TokenInfo[];
  error?: string;
}

interface ParseResponse {
  ast?: ParsedAST | null;
  error?: string;
}

interface CompileResponse {
  instructions?: string;
  constants?: number;
  error?: string;
}

interface ExecuteResponse {
  result?: string;
  output?: string;
  error?: string;
}

declare global {
  interface Window {
    monkeyWasmReady?: boolean;
    monkeyTokenize?: (code: string) => any;
    monkeyParseAST?: (code: string) => any;
    monkeyCompile?: (code: string) => any;
    monkeyExecute?: (code: string) => any;
    monkeyRepl?: (code: string) => any;
    monkeyCleanup?: () => void;
    Go?: any;
  }
}

class WasmService {
  private wasmReady = false;
  private wasmPromise: Promise<void> | null = null;

  constructor() {
    this.initWasm();
  }

  private async initWasm(): Promise<void> {
    if (this.wasmPromise) {
      return this.wasmPromise;
    }

    this.wasmPromise = new Promise(async (resolve, reject) => {
      try {
        // Check if already loaded
        if (window.monkeyWasmReady) {
          this.wasmReady = true;
          resolve();
          return;
        }

        // Load the WASM exec script if not already loaded
        if (!window.Go) {
          await this.loadWasmScript();
        }

        // Initialize Go WASM runtime
        const go = new window.Go();

        // Load and instantiate the WASM module
        const wasmModule = await WebAssembly.instantiateStreaming(
          fetch("/monkey.wasm"),
          go.importObject
        );

        // Run the Go program
        go.run(wasmModule.instance);

        // Wait for WASM to be ready with timeout
        let attempts = 0;
        const maxAttempts = 100; // 10 seconds timeout

        const checkReady = () => {
          attempts++;
          if (window.monkeyWasmReady) {
            this.wasmReady = true;
            console.log("WASM initialized successfully");
            resolve();
          } else if (attempts >= maxAttempts) {
            reject(new Error("WASM initialization timeout"));
          } else {
            setTimeout(checkReady, 100);
          }
        };

        checkReady();
      } catch (error) {
        console.error("WASM initialization error:", error);
        reject(error);
      }
    });

    return this.wasmPromise;
  }

  private loadWasmScript(): Promise<void> {
    return new Promise((resolve, reject) => {
      const script = document.createElement("script");
      script.src = "/wasm_exec.js";
      script.onload = () => {
        console.log("WASM exec script loaded");
        resolve();
      };
      script.onerror = (error) => {
        console.error("Failed to load wasm_exec.js:", error);
        reject(error);
      };
      document.head.appendChild(script);
    });
  }

  private async ensureReady(): Promise<void> {
    if (!this.wasmReady) {
      await this.initWasm();
    }
  }

  async tokenize(code: string): Promise<TokenizeResponse> {
    await this.ensureReady();

    if (!window.monkeyTokenize) {
      return { error: "WASM tokenize function not available" };
    }

    try {
      const result = window.monkeyTokenize(code);
      return result;
    } catch (error) {
      console.error("Tokenization error:", error);
      return { error: `Tokenization error: ${error}` };
    }
  }

  async parse(code: string): Promise<ParseResponse> {
    await this.ensureReady();

    if (!window.monkeyParseAST) {
      return { error: "WASM parse function not available" };
    }

    try {
      const result = window.monkeyParseAST(code);
      return result;
    } catch (error) {
      console.error("Parse error:", error);
      return { error: `Parse error: ${error}` };
    }
  }

  async compile(code: string): Promise<CompileResponse> {
    await this.ensureReady();

    if (!window.monkeyCompile) {
      return { error: "WASM compile function not available" };
    }

    try {
      const result = window.monkeyCompile(code);
      return result;
    } catch (error) {
      console.error("Compile error:", error);
      return { error: `Compile error: ${error}` };
    }
  }

  async execute(code: string): Promise<ExecuteResponse> {
    await this.ensureReady();

    console.log("WASM ready state:", this.wasmReady);
    console.log("window.monkeyWasmReady:", window.monkeyWasmReady);
    console.log("window.monkeyExecute exists:", !!window.monkeyExecute);
    console.log("typeof window.monkeyExecute:", typeof window.monkeyExecute);

    if (!window.monkeyExecute) {
      return { error: "WASM execute function not available" };
    }

    try {
      const result = window.monkeyExecute(code);
      console.log("WASM execute result:", result);
      console.log("WASM execute result type:", typeof result);
      console.log("WASM execute result is null:", result === null);
      console.log("WASM execute result is undefined:", result === undefined);

      if (!result || typeof result !== "object") {
        return {
          error: `WASM execute returned invalid response: ${typeof result}, value: ${result}`,
        };
      }

      return result;
    } catch (error) {
      console.error("Execution error:", error);
      return { error: `Execution error: ${error}` };
    }
  }

  async repl(code: string): Promise<ExecuteResponse> {
    await this.ensureReady();

    if (!window.monkeyRepl) {
      return { error: "WASM repl function not available" };
    }

    try {
      const result = window.monkeyRepl(code);
      console.log("WASM repl result:", result);

      if (!result || typeof result !== "object") {
        return { error: "WASM repl returned invalid response" };
      }

      return result;
    } catch (error) {
      console.error("REPL error:", error);
      return { error: `REPL error: ${error}` };
    }
  }
}

// Export singleton instance
export const wasmService = new WasmService();
export type {
  TokenInfo,
  ParsedAST,
  TokenizeResponse,
  ParseResponse,
  CompileResponse,
  ExecuteResponse,
};
