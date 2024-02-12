package parser

import (
	"eggo/token"
	"fmt"
)

type ASTnode struct {
	Data       token.Token
	Left       *ASTnode
	Right      *ASTnode
	IsTerminal bool
}

// TODO find a clean way to print the recursive struct definition
func (a *ASTnode) String() string {
	return fmt.Sprintf("{Data: %s Left:%v Right:%v}", a.Data.Literal, a.Left, a.Right)
}
