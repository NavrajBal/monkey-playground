package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey-wasm/compiler"
	"monkey-wasm/lexer"
	"monkey-wasm/object"
	"monkey-wasm/parser"
	"monkey-wasm/vm"
)

const PROMPT = ">> "

// Start runs a REPL backed by compiler + VM with persistent state
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)

	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins { symbolTable.DefineBuiltin(i, v.Name) }

	for {
		fmt.Fprint(out, PROMPT)
		if !scanner.Scan() { return }

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 { printParserErrors(out, p.Errors()); continue }

		comp := compiler.NewWithState(symbolTable, constants)
		if err := comp.Compile(program); err != nil { fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err); continue }

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		if err := machine.Run(); err != nil { fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err); continue }

		last := machine.LastPoppedStackElem()
		io.WriteString(out, last.Inspect())
		io.WriteString(out, "\n")
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-\''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors { io.WriteString(out, "\t"+msg+"\n") }
}


