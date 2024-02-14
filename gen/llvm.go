package gen

import (
	"eggo/ast"
	"eggo/parser"
	"eggo/token"
	"fmt"
	"log"
	"os"
)

const (
	preamble = `; ModuleID = 'examples/test1'
source_filename = "examples/test1"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-pc-linux-gnu"

@print_int_fstring = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
`

	postamble = `	ret i32 0
}

declare i32 @printf(i8*, ...) #1

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"PIC Level", i32 2}
!2 = !{i32 7, !"PIE Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 1}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{!"Ubuntu clang version 10.0.0-4ubuntu1"}
`
)

var operatorTable = map[string]string{
	token.PLUS:  "add",
	token.MINUS: "sub",
	token.STAR:  "mul",
	token.SLASH: "div",
}
var numRegisters = 0
var stackPosition = 0
var llvm_gen = ""

func GenerateLLVM(inFilepath string, outfilepath string) {
	p := parser.New(inFilepath)

	llvm_gen += preamble
	for root := p.ParseStatement(); root != nil; root = p.ParseStatement() {
		// each statement has its own AST
		gen_StackAllocations(root)
		gen_AST(root)
	}
	llvm_gen += postamble

	err := os.WriteFile(outfilepath, []byte(llvm_gen), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully generated llvm!\n")
}

func gen_StackAllocations(root *ast.ASTnode) {
	var worker func(node *ast.ASTnode)
	worker = func(node *ast.ASTnode) {
		if node == nil {
			return
		}

		worker(node.Left)
		worker(node.Right)

		if node.IsTerminal {
			numRegisters += 1
			llvm_gen += fmt.Sprintf("\t%%%d = alloca i32\n", numRegisters)
		}
	}
	worker(root)
	stackPosition = numRegisters
}

func gen_AST(root *ast.ASTnode) {
	nodeIsOperator := func(node *ast.ASTnode) bool {
		return node.Token.Type == token.PLUS || node.Token.Type == token.MINUS || node.Token.Type == token.STAR || node.Token.Type == token.SLASH
	}

	switch root.Token.Type {
	case token.PRINT:
		gen_print(root)
	default:
		if nodeIsOperator(root) {
			gen_expression(root)
		} else {
			log.Fatalf("Unexpected ast token type in llvm %s\n", root.Token.Type)
		}
	}
}

func gen_print(node *ast.ASTnode) int {
	argReg := gen_expression(node.Left)

	numRegisters += 1
	llvm_gen += fmt.Sprintf("\t%%%d = call i32(i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @print_int_fstring, i32 0, i32 0), i32 %%%d)\n", numRegisters, argReg)
	return numRegisters
}

func gen_expression(root *ast.ASTnode) int {
	if root.Token.Type == token.INT {
		llvm_gen += fmt.Sprintf("\tstore i32 %s, i32* %%%d\n", root.Token.Literal, stackPosition)
		stackPosition -= 1
		return stackPosition + 1
	}

	leftReg := gen_expression(root.Left)
	rightReg := gen_expression(root.Right)

	if root.Left.IsTerminal {
		numRegisters += 1
		llvm_gen += fmt.Sprintf("\t%%%d = load i32, i32* %%%d\n", numRegisters, leftReg)
		leftReg = numRegisters
	}

	if root.Right.IsTerminal {
		numRegisters += 1
		llvm_gen += fmt.Sprintf("\t%%%d = load i32, i32* %%%d\n", numRegisters, rightReg)
		rightReg = numRegisters
	}

	numRegisters += 1
	llvm_gen += fmt.Sprintf("\t%%%d = %s nsw i32 %%%d, %%%d\n", numRegisters, operatorTable[string(root.Token.Type)], leftReg, rightReg)

	return numRegisters
}
