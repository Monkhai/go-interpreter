package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkPraserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got= %d", len(program.Statements))
	}

	tests := []struct{ expectedIdentifier string }{
		{expectedIdentifier: "x"},
		{expectedIdentifier: "y"},
		{expectedIdentifier: "foobar"},
	}

	for i, tt := range tests {
		s := program.Statements[i]

		if !testLetStatement(t, s, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkPraserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, s := range program.Statements {
		returnStatement, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("statement not *ast.returnStatement. got=%T", s)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("retrunStatement.TokenLiteral() not 'return'. got=%q", returnStatement.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%T", s)
		return false
	}
	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value is not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() is not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func checkPraserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
