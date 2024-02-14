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

// func (a *ASTnode) String() string {
// 	var formatedString func(node *ASTnode, level int) string
// 	formatedString = func(node *ASTnode, level int) string {
// 		if node == nil {
// 			return "nil"
// 		}

// 		prev_indent := strings.Repeat("  ", level)
// 		cur_indent := strings.Repeat("  ", level+1)

// 		lBrace := "{\n"
// 		token := fmt.Sprintf("%s\"Token\": \"%s\"\n", cur_indent, node.Token.Literal)
// 		isTerminal := fmt.Sprintf("%s\"IsTerminal\": \"%t\"\n", cur_indent, node.IsTerminal)
// 		left := fmt.Sprintf("%s\"Left\": \"%s\"\n", cur_indent, formatedString(node.Left, level+1))
// 		right := fmt.Sprintf("%s\"Right\": \"%s\"\n", cur_indent, formatedString(node.Right, level+1))
// 		rBrace := fmt.Sprintf("%s}", prev_indent)

// 		return lBrace + token + isTerminal + left + right + rBrace
// 	}
// 	return formatedString(a, 0)
// }

func (a *ASTnode) String() string {
	var formatedString func(node *ASTnode, level int) string
	formatedString = func(node *ASTnode, level int) string {
		if node == nil {
			return "nil"
		}

		prev_indent := strings.Repeat("  ", level)
		cur_indent := strings.Repeat("  ", level+1)

		lBrace := "{\n"
		token := fmt.Sprintf("%sToken: %s\n", cur_indent, node.Token.Literal)
		isTerminal := fmt.Sprintf("%sIsTerminal: %t\n", cur_indent, node.IsTerminal)
		left := fmt.Sprintf("%sLeft: %s\n", cur_indent, formatedString(node.Left, level+1))
		right := fmt.Sprintf("%sRight: %s\n", cur_indent, formatedString(node.Right, level+1))
		rBrace := fmt.Sprintf("%s}", prev_indent)

		return lBrace + token + isTerminal + left + right + rBrace
	}
	return formatedString(a, 0)
}
