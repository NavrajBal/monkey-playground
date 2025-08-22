import axios from "axios";
import { type ParsedAST } from "../types/ast";

const API_BASE_URL = "http://localhost:8080/api";

export interface TokenInfo {
  type: string;
  literal: string;
  position: number;
}

export interface ExecuteResponse {
  result: string;
  output?: string;
  error?: string;
}

export interface TokenizeResponse {
  tokens: TokenInfo[];
  error?: string;
}

export interface ParseResponse {
  ast: ParsedAST | null;
  error?: string;
}

export interface CompileResponse {
  bytecode: number[];
  constants: string[];
  instructions: string;
  error?: string;
}

class ApiService {
  async execute(code: string): Promise<ExecuteResponse> {
    try {
      const response = await axios.post(`${API_BASE_URL}/execute`, { code });
      return response.data;
    } catch (error) {
      console.error("Execute error:", error);
      return { result: "", error: "Failed to execute code" };
    }
  }

  async tokenize(code: string): Promise<TokenizeResponse> {
    try {
      const response = await axios.post(`${API_BASE_URL}/tokenize`, { code });
      return response.data;
    } catch (error) {
      console.error("Tokenize error:", error);
      return { tokens: [], error: "Failed to tokenize code" };
    }
  }

  async parse(code: string): Promise<ParseResponse> {
    try {
      const response = await axios.post(`${API_BASE_URL}/parse`, { code });
      return response.data;
    } catch (error) {
      console.error("Parse error:", error);
      return { ast: null, error: "Failed to parse code" };
    }
  }

  async compile(code: string): Promise<CompileResponse> {
    try {
      const response = await axios.post(`${API_BASE_URL}/compile`, { code });
      return response.data;
    } catch (error) {
      console.error("Compile error:", error);
      return {
        bytecode: [],
        constants: [],
        instructions: "",
        error: "Failed to compile code",
      };
    }
  }

  async repl(code: string): Promise<ExecuteResponse> {
    try {
      const response = await axios.post(`${API_BASE_URL}/repl`, { code });
      return response.data;
    } catch (error) {
      console.error("REPL error:", error);
      return { result: "", error: "Failed to execute in REPL" };
    }
  }
}

export const apiService = new ApiService();
