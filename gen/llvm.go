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

var dataTypeTable = map[string]string{
	"bool":  "i1",
	"char":  "i8",
	"short": "i16",
	"int":   "i32",
	"long":  "i64",
}

type symbolTableType map[string]*Symbol

var globalSymbolTable symbolTableType

type llvmRegister struct {
	// thinking about register number
	// data type
	// maybe a symbol can have an associated register
}

// var stackPosition = 0
var numRegisters = 0
var llvm_gen = ""
var p parser.Parser

func GenerateLLVM(inFilepath string, outfilepath string) {
	p = *parser.New(inFilepath)

	llvm_gen += preamble
	gen_statements(p.ParseStatement(), globalSymbolTable)
	llvm_gen += postamble

	err := os.WriteFile(outfilepath, []byte(llvm_gen), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully generated llvm!\n")
}

func gen_statements(root ast.ASTnode, symbolTable symbolTableType) string {
	for ; root != nil; root = p.ParseStatement() {
		gen_statement(root, symbolTable)
	}
}

func gen_statement(root ast.ASTnode, symbolTable symbolTableType) string {
	switch root := root.(type) {
	case *ast.DeclareNode:
		return gen_declaration(root, symbolTable)
	case *ast.AssignNode:
		return gen_assign(root, symbolTable)
	case *ast.PrintNode:
		return gen_print(root, symbolTable)
	case *ast.BlockNode:
		return gen_block(root, symbolTable)
	case *ast.IfNode:
		return gen_if(root, symbolTable)
	case *ast.WhileNode:
		return gen_while(root, symbolTable)
	default:
		fmt.Printf("Unexpeded statement type!\n")
	}
}

func gen_declaration(root *ast.DeclareNode, symbolTable symbolTableType) string {
	symbolName, dataType := root.Ident, dataTypeTable[root.DataType]
	symbolTable[root.Ident] = &Symbol{Name: symbolName, DataType: dataType}
	llvm_gen += fmt.Sprintf("\t%%%s = alloca %s\n", symbolName, dataType)

	return fmt.Sprintf("%%%s", symbolName)
}

func gen_assign(root *ast.AssignNode, symbolTable symbolTableType) {
	exprValue := gen_expression(root.Expression, symbolTable)
	exprDataType := "i32" // for now just let all expr data types be i32
	Smbl := symbolTable[root.Ident]
	if Smbl == nil {
		log.Fatalf("LLVM Error Symbol \"%s\" not found!\n", root.Ident)
	}
	llvm_gen += fmt.Sprintf("\tstore %s %s, %s* %%%s\n", exprDataType, exprValue, Smbl.DataType, Smbl.Name)
}

func gen_print(node *ast.PrintNode, symbolTable symbolTableType) string {
	argReg := gen_expression(node.Expression, symbolTable)

	numRegisters += 1
	llvm_gen += fmt.Sprintf("\t%%%d = call i32(i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @print_int_fstring, i32 0, i32 0), i32 %s)\n", numRegisters, argReg)
	return fmt.Sprintf("%%%d", numRegisters)
}

func gen_expression(root *ast.ExpressionNode, symbolTable symbolTableType) string {
	if root.IsTerminal() {
		if Smbl := symbolTable[root.Value]; Smbl != nil {
			numRegisters += 1
			llvm_gen += fmt.Sprintf("\t%%%d = load %s, %s* %%%s\n", numRegisters, Smbl.DataType, Smbl.DataType, Smbl.Name)
			return fmt.Sprintf("%%%d", numRegisters)
		} else if _, err := strconv.Atoi(root.Value); err == nil {
			return string(root.Value)
		} else {
			log.Fatalf("Invalid Symbol in expression: %s\n", root.Value)
		}
	}

	type OpExprPair struct {
		Op       string
		ExprType string
	}
	OpExprTable := map[string]OpExprPair{
		token.PLUS: {
			Op:       "add",
			ExprType: "numerical",
		},
		token.MINUS: {
			Op:       "sub",
			ExprType: "numerical",
		},
		token.STAR: {
			Op:       "mul",
			ExprType: "numerical",
		},
		token.SLASH: {
			Op:       "div",
			ExprType: "numerical",
		},
		// TODO - figure out how to add shifting in LLVM
		// token.LSHIFT: {
		// 	Op:    "lshl",
		// 	ExprType: "numerical",
		// },
		// token.RSHIFT: {
		// 	Op:    "lshr",
		// 	ExprType: "numerical",
		// },
		token.EQ: {
			Op:       "eq",
			ExprType: "boolean",
		},
		token.NE: {
			Op:       "ne",
			ExprType: "boolean",
		},
		token.LT: {
			Op:       "slt",
			ExprType: "boolean",
		},
		token.LE: {
			Op:       "sle",
			ExprType: "boolean",
		},
		token.GT: {
			Op:       "sgt",
			ExprType: "boolean",
		},
		token.GE: {
			Op:       "sge",
			ExprType: "boolean",
		},
	}
	operator, expressionType := OpExprTable[root.Value].Op, OpExprTable[root.Value].ExprType

	leftReg := gen_expression(root.Left, symbolTable)
	rightReg := gen_expression(root.Right, symbolTable)

	numRegisters += 1
	switch expressionType {
	case "numerical":
		llvm_gen += fmt.Sprintf("\t%%%d = %s nsw i32 %s, %s\n", numRegisters, operator, leftReg, rightReg)
	case "boolean":
		llvm_gen += fmt.Sprintf("\t%%%d = icmp %s i32 %s, %s\n", numRegisters, operator, leftReg, rightReg)
		// sign extend the bool
		numRegisters += 1
		llvm_gen += fmt.Sprintf("\t%%%d = zext i1 %d to i32\n", numRegisters, numRegisters-1)
	default:
		log.Fatalf("Unexpeded expression type %s", expressionType)
	}

	return fmt.Sprintf("%%%d", numRegisters)
}

func gen_block(root *ast.BlockNode, smblTable symbolTableType) string {
	return gen_statements(root.Statements, smblTable)
}

func gen_if(root *ast.IfNode, smblTable symbolTableType) {
	condition := gen_expression(root.Condition, smblTable)
	fmt.Sprintf("\tbr i1 %s, label %%true, label %%false\n", condition)
	fmt.Sprintf("true:\n")
	gen_block(root.HappyBody, smblTable)
	fmt.Sprintf("br label %%tail\n")
	if root.ContainsElse() {
		fmt.Sprintf("false:\n")
		fmt.Sprintf("br label %%tail\n")
	}
	fmt.Sprintf("tail:\n")
}

func gen_while(root *ast.WhileNode, smblTable symbolTableType) {
	condition := gen_expression(root.Condition, smblTable)
	fmt.Sprintf("br label %%condition\n")
	fmt.Sprintf("condition:\n")
	fmt.Sprintf("br i1 %%%s, label %%body, label %%tail", condition)
	fmt.Sprintf("body:\n")
	gen_block(root.Body, smblTable)
	fmt.Sprintf("br label, %%condition\n")
	fmt.Sprintf("tail:\n")
}
