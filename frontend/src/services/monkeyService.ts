import { config, isUsingWasm } from "../config/config";
import { apiService } from "./api";
import { wasmService } from "./wasmService";
import type { TokenInfo } from "./wasmService";

// Unified interfaces that work with both backends
export interface TokenizeResponse {
  tokens?: TokenInfo[];
  error?: string;
}

export interface ParseResponse {
  ast?: any;
  error?: string;
}

export interface CompileResponse {
  instructions?: string;
  constants?: number | string[]; // Support both formats
  bytecode?: number[];
  error?: string;
}

export interface ExecuteResponse {
  result?: string;
  output?: string;
  error?: string;
}

/**
 * Unified service that can switch between API and WASM backends
 */
class MonkeyService {
  constructor() {
    // Use the existing apiService instance
  }

  async tokenize(code: string): Promise<TokenizeResponse> {
    if (isUsingWasm()) {
      return wasmService.tokenize(code);
    } else {
      const result = await apiService.tokenize(code);
      return {
        tokens: result.tokens,
        error: result.error,
      };
    }
  }

  async parse(code: string): Promise<ParseResponse> {
    if (isUsingWasm()) {
      return wasmService.parse(code);
    } else {
      const result = await apiService.parse(code);
      return {
        ast: result.ast,
        error: result.error,
      };
    }
  }

  async compile(code: string): Promise<CompileResponse> {
    if (isUsingWasm()) {
      return wasmService.compile(code);
    } else {
      const result = await apiService.compile(code);
      return {
        instructions: result.instructions,
        constants: result.constants,
        bytecode: result.bytecode,
        error: result.error,
      };
    }
  }

  async execute(code: string): Promise<ExecuteResponse> {
    if (isUsingWasm()) {
      return wasmService.execute(code);
    } else {
      const result = await apiService.execute(code);
      return {
        result: result.result,
        output: result.output,
        error: result.error,
      };
    }
  }

  async repl(code: string): Promise<ExecuteResponse> {
    if (isUsingWasm()) {
      return wasmService.repl(code);
    } else {
      const result = await apiService.repl(code);
      return {
        result: result.result,
        output: result.output,
        error: result.error,
      };
    }
  }

  // Get current backend info
  getBackendInfo() {
    return {
      backend: config.backend,
      isWasm: isUsingWasm(),
      apiUrl: config.apiUrl,
    };
  }
}

// Export singleton instance
export const monkeyService = new MonkeyService();
export type { TokenInfo };
