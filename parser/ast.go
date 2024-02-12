package parser

import (
	"eggo/token"
	"fmt"
	"strings"
)

type ASTnode struct {
	Token      token.Token
	Left       *ASTnode
	Right      *ASTnode
	IsTerminal bool
}

func (a *ASTnode) String() string {
	return a.stringWithIndent(0) + "\n"
}

func (a *ASTnode) stringWithIndent(level int) string {
	if a == nil {
		return "nil"
	}

	prev_indent := strings.Repeat("  ", level)
	cur_indent := strings.Repeat("  ", level+1)

	lBrace := "{\n"
	token := fmt.Sprintf("%sToken: %s\n", cur_indent, a.Token.Literal)
	isTerminal := fmt.Sprintf("%sIsTerminal: %t\n", cur_indent, a.IsTerminal)
	left := fmt.Sprintf("%sLeft: %v\n", cur_indent, a.Left.stringWithIndent(level+1))
	right := fmt.Sprintf("%sRight: %v\n", cur_indent, a.Right.stringWithIndent(level+1))
	rBrace := fmt.Sprintf("%s}", prev_indent)

	return lBrace + token + isTerminal + left + right + rBrace
}
