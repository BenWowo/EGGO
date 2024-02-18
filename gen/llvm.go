package gen

import (
	"eggo/ast"
	"eggo/parser"
	"eggo/token"
	"fmt"
	"log"
	"os"
	"strconv"
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

type Symbol struct {
	Name     string
	DataType string
}

var SymbolTable = map[string]*Symbol{}

// var stackPosition = 0
var numRegisters = 0
var llvm_gen = ""

func GenerateLLVM(inFilepath string, outfilepath string) {
	p := parser.New(inFilepath)

	llvm_gen += preamble
	for root := p.ParseStatement(); root != nil; root = p.ParseStatement() {
		switch root := root.(type) {
		case *ast.DeclareNode:
			gen_declaration(root)
		case *ast.AssignNode:
			gen_assign(root)
		case *ast.PrintNode:
			gen_print(root)
		default:
			fmt.Printf("Unexpeded statement type!\n")
		}
	}
	llvm_gen += postamble

	err := os.WriteFile(outfilepath, []byte(llvm_gen), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully generated llvm!\n")
}

func gen_declaration(root *ast.DeclareNode) string {
	DataTypeTable := map[string]string{
		"int": "i32",
	}

	SymbolName, DataType := root.Ident, DataTypeTable[root.DataType]
	SymbolTable[root.Ident] = &Symbol{Name: SymbolName, DataType: DataType}
	llvm_gen += fmt.Sprintf("\t%%%s = alloca %s\n", SymbolName, DataType)

	return fmt.Sprintf("%%%s", SymbolName)
}

func gen_assign(root *ast.AssignNode) {
	exprValue := gen_expression(root.Expression)
	exprDataType := "i32" // for now just let all expr data types be i32
	Smbl := SymbolTable[root.Ident]
	if Smbl == nil {
		log.Fatalf("LLVM Error Symbol \"%s\" not found!\n", root.Ident)
	}
	llvm_gen += fmt.Sprintf("\tstore %s %s, %s* %%%s\n", exprDataType, exprValue, Smbl.DataType, Smbl.Name)
}

func gen_print(node *ast.PrintNode) string {
	argReg := gen_expression(node.Expression)

	numRegisters += 1
	llvm_gen += fmt.Sprintf("\t%%%d = call i32(i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @print_int_fstring, i32 0, i32 0), i32 %s)\n", numRegisters, argReg)
	return fmt.Sprintf("%%%d", numRegisters)
}

func gen_expression(root *ast.ExpressionNode) string {
	if root.IsTerminal() {
		Smbl := SymbolTable[root.Value]
		if Smbl != nil {
			numRegisters += 1
			llvm_gen += fmt.Sprintf("\t%%%d = load %s, %s* %%%s\n", numRegisters, Smbl.DataType, Smbl.DataType, Smbl.Name)
			return fmt.Sprintf("%%%d", numRegisters)
		} else if _, err := strconv.Atoi(root.Value); err == nil {
			return string(root.Value)
		} else {
			log.Fatalf("Invalid Symbol in expression: %s\n", root.Value)
		}
	}

	leftReg := gen_expression(root.Left)
	rightReg := gen_expression(root.Right)

	OperatorTable := map[string]string{
		token.PLUS:  "add",
		token.MINUS: "sub",
		token.STAR:  "mul",
		token.SLASH: "div",
		// token.LSHIFT: "lshl",
		// token.RSHIFT: "lshr",
	}

	numRegisters += 1
	llvm_gen += fmt.Sprintf("\t%%%d = %s nsw i32 %s, %s\n", numRegisters, OperatorTable[root.Value], leftReg, rightReg)

	return fmt.Sprintf("%%%d", numRegisters)
}
