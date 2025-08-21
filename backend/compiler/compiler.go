package compiler

import (
	"fmt"
	"monkey-playground-backend/ast"
	"monkey-playground-backend/code"
	"monkey-playground-backend/object"
	"sort"
)

type Compiler struct {
	constants []object.Object

	scopes     []CompilationScope
	scopeIndex int

	symbolTable *SymbolTable
}

func New() *Compiler {
	mainScope := CompilationScope{instructions: code.Instructions{}}
	c := &Compiler{constants: []object.Object{}, scopes: []CompilationScope{mainScope}, scopeIndex: 0, symbolTable: NewSymbolTable()}
	// define builtins for compiler resolution
	for i, b := range object.Builtins { c.symbolTable.DefineBuiltin(i, b.Name) }
	return c
}

func NewWithState(table *SymbolTable, constants []object.Object) *Compiler {
	mainScope := CompilationScope{instructions: code.Instructions{}}
	c := &Compiler{constants: constants, scopes: []CompilationScope{mainScope}, scopeIndex: 0, symbolTable: table}
	return c
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements { if err := c.Compile(s); err != nil { return err } }

	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil { return err }
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		if node.Operator == "<" {
			if err := c.Compile(node.Right); err != nil { return err }
			if err := c.Compile(node.Left); err != nil { return err }
			c.emit(code.OpGreaterThan)
			return nil
		}
		if err := c.Compile(node.Left); err != nil { return err }
		if err := c.Compile(node.Right); err != nil { return err }
		switch node.Operator {
		case "+": c.emit(code.OpAdd)
		case "-": c.emit(code.OpSub)
		case "*": c.emit(code.OpMul)
		case "/": c.emit(code.OpDiv)
		case ">": c.emit(code.OpGreaterThan)
		case "==": c.emit(code.OpEqual)
		case "!=": c.emit(code.OpNotEqual)
		default: return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value { c.emit(code.OpTrue) } else { c.emit(code.OpFalse) }

	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil { return err }
		switch node.Operator {
		case "!": c.emit(code.OpBang)
		case "-": c.emit(code.OpMinus)
		default: return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		if err := c.Compile(node.Condition); err != nil { return err }
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)
		if err := c.Compile(node.Consequence); err != nil { return err }
		if c.lastInstructionIs(code.OpPop) { c.replaceLastPopWithReturn() }
		jumpPos := c.emit(code.OpJump, 9999)
		afterConsequence := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequence)
		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			if err := c.Compile(node.Alternative); err != nil { return err }
			if c.lastInstructionIs(code.OpPop) { c.replaceLastPopWithReturn() }
		}
		afterAlternative := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternative)

	case *ast.BlockStatement:
		for _, s := range node.Statements { if err := c.Compile(s); err != nil { return err } }

	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil { return err }
		symbol := c.symbolTable.Define(node.Name.Value)
		switch symbol.Scope {
		case GlobalScope:
			c.emit(code.OpSetGlobal, symbol.Index)
		case LocalScope:
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok { return fmt.Errorf("undefined variable %s", node.Value) }
		c.loadSymbol(symbol)

	case *ast.ReturnStatement:
		if err := c.Compile(node.ReturnValue); err != nil { return err }
		c.emit(code.OpReturnValue)

	case *ast.FunctionLiteral:
		c.enterScope()
		if node.Name != "" { c.symbolTable.DefineFunctionName(node.Name) }
		for _, p := range node.Parameters { c.symbolTable.Define(p.Value) }
		if err := c.Compile(node.Body); err != nil { return err }
		if c.lastInstructionIs(code.OpPop) { c.replaceLastPopWithReturn() }
		if !c.lastInstructionIs(code.OpReturnValue) { c.emit(code.OpReturn) }
		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		ins := c.leaveScope()
		for _, s := range freeSymbols { c.loadSymbol(s) }
		compiledFn := &object.CompiledFunction{Instructions: ins, NumLocals: numLocals, NumParameters: len(node.Parameters)}
		fnIndex := c.addConstant(compiledFn)
		c.emit(code.OpClosure, fnIndex, len(freeSymbols))

	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil { return err }
		for _, a := range node.Arguments { if err := c.Compile(a); err != nil { return err } }
		c.emit(code.OpCall, len(node.Arguments))

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements { if err := c.Compile(el); err != nil { return err } }
		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs { keys = append(keys, k) }
		sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
		for _, k := range keys {
			if err := c.Compile(k); err != nil { return err }
			if err := c.Compile(node.Pairs[k]); err != nil { return err }
		}
		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.IndexExpression:
		if err := c.Compile(node.Left); err != nil { return err }
		if err := c.Compile(node.Index); err != nil { return err }
		c.emit(code.OpIndex)
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode { return &Bytecode{Instructions: c.currentInstructions(), Constants: c.constants} }

func (c *Compiler) addConstant(obj object.Object) int { c.constants = append(c.constants, obj); return len(c.constants) - 1 }

func (c *Compiler) emit(op code.Opcode, operands ...int) int { ins := code.Make(op, operands...); pos := c.addInstruction(ins); c.setLastInstruction(op, pos); return pos }

func (c *Compiler) addInstruction(ins []byte) int { pos := len(c.currentInstructions()); c.scopes[c.scopeIndex].instructions = append(c.currentInstructions(), ins...); return pos }

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) { previous := c.scopes[c.scopeIndex].lastInstruction; last := EmittedInstruction{Opcode: op, Position: pos}; c.scopes[c.scopeIndex].previousInstruction = previous; c.scopes[c.scopeIndex].lastInstruction = last }

func (c *Compiler) lastInstructionIs(op code.Opcode) bool { return len(c.currentInstructions()) > 0 && c.scopes[c.scopeIndex].lastInstruction.Opcode == op }

func (c *Compiler) replaceLastPopWithReturn() { lastPos := c.scopes[c.scopeIndex].lastInstruction.Position; c.replaceInstruction(lastPos, code.Make(code.OpReturnValue)); c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue }

func (c *Compiler) removeLastPop() { c.scopes[c.scopeIndex].instructions = c.currentInstructions()[:c.scopes[c.scopeIndex].lastInstruction.Position]; c.scopes[c.scopeIndex].lastInstruction = c.scopes[c.scopeIndex].previousInstruction }

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) { for i := 0; i < len(newInstruction); i++ { c.scopes[c.scopeIndex].instructions[pos+i] = newInstruction[i] } }

func (c *Compiler) changeOperand(opPos int, operand int) { op := code.Opcode(c.currentInstructions()[opPos]); newInstruction := code.Make(op, operand); c.replaceInstruction(opPos, newInstruction) }

func (c *Compiler) enterScope() { scope := CompilationScope{instructions: code.Instructions{}}; c.scopes = append(c.scopes, scope); c.scopeIndex++; c.symbolTable = NewEnclosedSymbolTable(c.symbolTable) }

func (c *Compiler) leaveScope() code.Instructions { instructions := c.currentInstructions(); c.scopes = c.scopes[:len(c.scopes)-1]; c.scopeIndex--; c.symbolTable = c.symbolTable.Outer; return instructions }

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	case FreeScope:
		c.emit(code.OpGetFree, s.Index)
	case FunctionScope:
		c.emit(code.OpCurrentClosure)
	}
}

func (c *Compiler) currentInstructions() code.Instructions { return c.scopes[c.scopeIndex].instructions }

type Bytecode struct { Instructions code.Instructions; Constants []object.Object }

type EmittedInstruction struct { Opcode code.Opcode; Position int }

type CompilationScope struct { instructions code.Instructions; lastInstruction EmittedInstruction; previousInstruction EmittedInstruction }


