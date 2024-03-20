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

var dataTypeTable = map[string]string{
	"bool":  "i1",
	"char":  "i8",
	"short": "i16",
	"int":   "i32",
	"long":  "i64",
}

type Symbol struct {
	Name     string
	DataType string
}

type smblTbleType map[string]*Symbol

var globalSmblTble = make(smblTbleType)

type llvmReg struct {
	IsImmediate bool
	DataType    string
	// Name of llvm register or immediate.
	// Example Names include "1" or "%jim".
	Name string
}

// var stackPosition = 0
var numRegisters = 0
var llvm_gen = ""
var p parser.Parser

// Generates llvm and writes to specified output file.
func GenerateLLVM(inFilePath string, outFilePath string) {
	p = *parser.New(inFilePath)
	root := p.ParseStatement()
	llvm_gen = preamble + genStatements(root, globalSmblTble) + postamble

	// TODO - compare performance between buffered file write and string write
	err := os.WriteFile(outFilePath, []byte(llvm_gen), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully generated llvm!\n")
}

// Generates and returns llvm for all of the statements within
// the source file.
func genStatements(root *ast.ASTnode, smblTble smblTbleType) string {
	var stmts string

	for ; root != nil; root = p.ParseStatement() {
		stmts += genStatement(root, smblTble)
	}

	return stmts
}

// Generates and returns llvm for a generic ASTnode.
func genStatement(root *ast.ASTnode, smblTble smblTbleType) (stmt string) {
	switch root := (*root).(type) {
	case *ast.DeclareNode:
		stmt = genDeclaration(root, smblTble)
	case *ast.AssignNode:
		stmt, _ = genAssign(root, smblTble)
	case *ast.PrintNode:
		stmt = genPrint(root, smblTble)
	case *ast.BlockNode:
		stmt = genBlock(root, smblTble)
	case *ast.IfNode:
		stmt = genIf(root, smblTble)
	case *ast.WhileNode:
		stmt = genWhile(root, smblTble)
	default:
		log.Fatalf("LLVM unexpected stmt type: [%T]\n", root)
	}

	return stmt
}

// Generates and returns llvm for a declaration AST.
func genDeclaration(root *ast.DeclareNode, smblTble smblTbleType) string {
	smblName, smblDataType := root.Ident, dataTypeTable[root.DataType]
	smblTble[smblName] = &Symbol{
		Name:     smblName,
		DataType: smblDataType,
	}

	return fmt.Sprintf("\t%%%s = alloca %s\n", smblName, smblDataType)
}

// Generates and returns llvm for an assign AST.
func genAssign(root *ast.AssignNode, smblTble smblTbleType) (string, *llvmReg) {
	var assignStmt string

	exprGen, exprReg := genExpression(root.Expression, smblTble)
	assignStmt += exprGen

	Smbl, found := smblTble[root.Ident]
	if !found {
		log.Fatalf("LLVM error symbol [%s] not found\n", root.Ident)
	}

	assignStmt += fmt.Sprintf("\tstore %s %s, %s* %%%s\n",
		exprReg.DataType, exprReg.Name, Smbl.DataType, Smbl.Name,
	)

	return assignStmt, exprReg
}

// Generates and returns llvm for print AST.
func genPrint(node *ast.PrintNode, smblTble smblTbleType) string {
	var printStmt string

	argGen, argReg := genExpression(node.Expression, smblTble)
	printStmt += argGen

	numRegisters += 1
	printStmt += fmt.Sprintf("\t%%%d = call i32(i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @print_int_fstring, i32 0, i32 0), i32 %s)\n",
		numRegisters, argReg.Name)

	return printStmt
}

// Generates and returns llvm and llvmReg where result is stored.
func genExpression(root *ast.ExpressionNode, smblTble smblTbleType) (string, *llvmReg) {
	var exprStmt string
	exprReg := new(llvmReg)

	if root.IsTerminal() {
		if Smbl := smblTble[root.Value]; Smbl != nil {
			numRegisters += 1
			exprStmt += fmt.Sprintf("\t%%%d = load %s, %s* %%%s\n", numRegisters, Smbl.DataType,
				Smbl.DataType, Smbl.Name)

			exprReg.DataType = Smbl.DataType
			exprReg.IsImmediate = false
			exprReg.Name = fmt.Sprintf("%%%d", numRegisters)

			return exprStmt, exprReg
		} else if _, err := strconv.Atoi(root.Value); err == nil {
			exprReg.DataType = "i32"
			exprReg.IsImmediate = true
			exprReg.Name = string(root.Value)

			return exprStmt, exprReg
		} else {
			log.Fatalf("Invalid Symbol in expression: %s\n", root.Value)
		}
	}
	exprReg.IsImmediate = false

	leftGen, leftReg := genExpression(root.Left, smblTble)
	exprStmt += leftGen

	rightGen, rightReg := genExpression(root.Right, smblTble)
	exprStmt += rightGen

	operator := getExprOperator(root)
	exprType := getExprType(root)

	numRegisters += 1
	switch exprType {
	case "i32":
		exprReg.DataType = "i32"
		exprStmt += fmt.Sprintf("\t%%%d = %s nsw i32 %s, %s\n", numRegisters, operator,
			leftReg.Name, rightReg.Name)
	case "bool":
		exprReg.DataType = "bool"
		exprStmt += fmt.Sprintf("\t%%%d = icmp %s i32 %s, %s\n", numRegisters, operator,
			leftReg.Name, rightReg.Name)

		// sign extend the bool
		numRegisters += 1
		exprStmt += fmt.Sprintf("\t%%%d = zext i1 %d to i32\n", numRegisters, numRegisters-1)
	default:
		log.Fatalf("Unexpeded expression type [%s]", exprType)
	}

	exprReg.Name = fmt.Sprintf("%%%d", numRegisters)
	return exprStmt, exprReg
}

func getExprOperator(root *ast.ExpressionNode) (Operator string) {
	switch root.Value {
	case token.PLUS:
		Operator = "add"
	case token.MINUS:
		Operator = "sub"
	case token.STAR:
		Operator = "mul"
	case token.SLASH:
		Operator = "div"
	case token.EQ:
		Operator = "eq"
	case token.NE:
		Operator = "ne"
	case token.LT:
		Operator = "slt"
	case token.LE:
		Operator = "sle"
	case token.GT:
		Operator = "sgt"
	case token.GE:
		Operator = "sge"
	}

	return Operator
}

// TODO: For now the types will only be "i32" and "i1"
// I have some options how I can handle types here...
func getExprType(root *ast.ExpressionNode) (ExprType string) {
	switch root.Value {
	case token.PLUS:
		ExprType = "i32"
	case token.MINUS:
		ExprType = "i32"
	case token.STAR:
		ExprType = "i32"
	case token.SLASH:
		ExprType = "i32"
	case token.EQ:
		ExprType = "bool"
	case token.NE:
		ExprType = "bool"
	case token.LT:
		ExprType = "bool"
	case token.LE:
		ExprType = "bool"
	case token.GT:
		ExprType = "bool"
	case token.GE:
		ExprType = "bool"
	}

	return ExprType
}

// Generates llvm for block AST.
func genBlock(root *ast.BlockNode, smblTable smblTbleType) string {
	var blockStmts string

	for stmtAST := p.ParseStatement(); stmtAST != nil; stmtAST = p.ParseStatement() {
		blockStmts += genStatement(stmtAST, smblTable)
	}

	return blockStmts
}

// TODO: 08 Conditionals and Loops
func genIf(root *ast.IfNode, smblTable smblTbleType) string {
	var ifStmt string

	_, conditionReg := genExpression(root.Condition, smblTable)

	ifStmt += fmt.Sprintf("\tbr i1 %s, label %%true, label %%false\n", conditionReg.Name)
	ifStmt += "true:\n"
	ifStmt += genBlock(root.HappyBody, smblTable)
	ifStmt += fmt.Sprintf("br label %%tail\n")
	if root.ContainsElse() {
		ifStmt += "false:\n"
		ifStmt += genStatement(root.SadBody, smblTable)
		ifStmt += "br label %%tail\n"
	}
	ifStmt += "tail:\n"

	return ifStmt
}

// TODO: 08 Conditionals and Loops
func genWhile(root *ast.WhileNode, smblTable smblTbleType) string {
	var whileStmt string

	condition, _ := genExpression(root.Condition, smblTable)
	whileStmt += fmt.Sprintf("br label %%condition\n")
	whileStmt += "condition:\n"
	whileStmt += fmt.Sprintf("br i1 %%%s, label %%body, label %%tail", condition)
	whileStmt += "body:\n"
	whileStmt += genBlock(root.Body, smblTable)
	whileStmt += fmt.Sprintf("br label, %%condition\n")
	whileStmt += "tail:\n"

	return whileStmt
}
