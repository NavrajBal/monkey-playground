export interface ASTNode {
  type: string;
  [key: string]: unknown;
}

export interface ParsedAST {
  type: string;
  statements: ASTNode[];
  string: string;
}
