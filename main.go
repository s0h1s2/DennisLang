package main

import (
	"fmt"
	"io"
	"os"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/lexer"
	"github.com/s0h1s2/parser"
	"github.com/s0h1s2/resolver"
)

// import (
//
//	"fmt"
//	"os"
//	"reflect"
//	"strconv"
//
// )
//
// type Node interface{}
//
//	type Expr interface {
//		Node
//		exprNode()
//	}
//
//	type Stmt interface {
//		Node
//		stmtNode()
//	}
//
//	type ExprInt struct {
//		Value uint64
//	}
//
//	type ExprIdent struct {
//		Name string
//	}
//
//	type ExprBin struct {
//		Left  Expr
//		Right Expr
//		Op    uint8
//	}
//
//	type StmtLet struct {
//		Name     string
//		TypeName string
//		Init     Expr
//	}
//
// func (t *ExprInt) exprNode()   {}
// func (t *ExprBin) exprNode()   {}
// func (t *ExprIdent) exprNode() {}
//
//	func alignTo(size uint64, align uint64) uint64 {
//		return (size + align - 1) / align * align
//	}
//
//	func resolver(env *Environment, node Node) {
//		switch t := node.(type) {
//		case *ExprInt:
//			{
//				println(t.Value)
//			}
//		case *ExprIdent:
//			{
//				if !env.Lookup(t.Name) {
//					fmt.Printf("Variable '%s' not found.\n", t.Name)
//					panic("")
//				}
//				result := env.Get(t.Name).(*EntryVar)
//				println("variable=", t.Name, "size", result.AddrOffset)
//			}
//		case *ExprBin:
//			{
//
//				resolver(env, t.Left)
//				resolver(env, t.Right)
//			}
//		case *StmtLet:
//			{
//				if env.Lookup(t.Name) {
//					fmt.Printf("Can't redeclare variable '%s' more than once.", t.Name)
//					panic("")
//				}
//				if !env.Lookup(t.TypeName) {
//					fmt.Printf("Type '%s' couldn't be found.", t.TypeName)
//					panic("")
//				}
//				typ := &env.GetType(t.TypeName).Typee
//				env.AddrOffset += typ.Size
//				env.AddrOffset = alignTo(env.AddrOffset, typ.Alignment)
//				env.Add(t.Name, &EntryVar{AddrOffset: env.AddrOffset, Typee: typ})
//			}
//		default:
//			{
//				panic("Unreachable")
//			}
//		}
//	}
//
//	type Type struct {
//		Size      uint64
//		Alignment uint64
//	}
//
//	type Entry interface {
//		entry()
//	}
//
//	type EntryVar struct {
//		AddrOffset uint64
//		Typee      *Type
//	}
//
//	type EntryType struct {
//		Typee Type
//	}
//
// func (e *EntryVar) entry()  {}
// func (e *EntryType) entry() {}
//
//	type Environment struct {
//		env        map[string]Entry
//		AddrOffset uint64
//	}
//
//	func (env *Environment) registerType(name string, entry Entry) {
//		_, ok := entry.(*EntryType)
//		if !ok {
//			panic("invalid type in registerType")
//		}
//		env.env[name] = entry
//	}
//
//	func (env *Environment) Lookup(name string) bool {
//		if env.env[name] != nil {
//			return true
//		}
//		return false
//	}
//
//	func (env *Environment) Add(name string, entry Entry) {
//		// TODO: maybe check using Lookup method but for now only add entry to the hashmap.
//		env.env[name] = entry
//	}
//
//	func (env *Environment) GetType(name string) *EntryType {
//		return env.env[name].(*EntryType)
//	}
//
//	func (env *Environment) GetEntryVar(name string) *EntryVar {
//		return env.env[name].(*EntryVar)
//	}
//
//	func (env *Environment) Get(name string) Entry {
//		return env.env[name]
//	}
//
//	func codegen_node(node Node, env *Environment) string {
//		switch n := node.(type) {
//		case *StmtLet:
//			{
//				initValue := codegen_node(n.Init, env)
//				return fmt.Sprintf("mov [rbp-%d],%s", env.GetEntryVar(n.Name).AddrOffset, initValue)
//			}
//		case *ExprInt:
//			{
//				t := strconv.Itoa(int(n.Value))
//				return t
//			}
//		case *ExprIdent:
//			{
//				addr := env.GetEntryVar(n.Name)
//
//				typeSpecifier := ""
//				switch addr.Typee.Size {
//				case 1:
//					{
//						typeSpecifier = " byte "
//					}
//				case 2:
//					{
//						typeSpecifier = " word "
//					}
//				case 4:
//					{
//						typeSpecifier = " dword "
//					}
//				default:
//					{
//						typeSpecifier = " qword "
//					}
//				}
//				return fmt.Sprintf("%s[rbp-%d]", typeSpecifier, addr.AddrOffset)
//			}
//		case *ExprBin:
//			{
//				left := codegen_node(n.Left, env)
//				right := codegen_node(n.Left, env)
//				var opcode string
//				if n.Op == '+' {
//					opcode = "add rax,rcx"
//				}
//				if n.Op == '*' {
//					opcode = "mul rax,rcx"
//				}
//				return fmt.Sprintf("mov rax,%s\nmov rcx,%s\n%s", left, right, opcode)
//			}
//		default:
//			{
//				println(reflect.TypeOf(n).Name())
//				println(reflect.TypeOf(n).Size())
//				println(n)
//				return "nop"
//			}
//		}
//
// }
//
//	func codegen(env *Environment, program []Node) {
//		f, err := os.Create("output.S")
//		if err != nil {
//			panic("Unable to open file.")
//		}
//		fmt.Fprintln(f, "format ELF64 executable 3")
//		fmt.Fprintln(f, "segment readable executable")
//		fmt.Fprintln(f, "start:")
//		fmt.Fprintln(f, "push rbp")
//		fmt.Fprintln(f, "mov rbp,rsp")
//		for i := 0; i < len(program); i++ {
//			node := program[i]
//			fmt.Fprintln(f, codegen_node(node, env))
//		}
//		fmt.Fprintln(f, "mov rax,1")
//		fmt.Fprintln(f, "int 0x80")
//	}
func printAst(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.ExprAssign:
		{
			printAst(e.Left)
			print("=")
			printAst(e.Right)
		}
	case *ast.ExprInt:
		{
			println(e.Value)
		}

	case *ast.ExprBinary:
		{
			printAst(e.Left)
			printAst(e.Right)
		}

	case *ast.ExprIdent:
		{
			println(e.Name)
		}
	default:
		{
			panic("UnReachable")
		}
	}
}
func hello(n ast.Node) bool {
	println(n)
	return false
}
func main() {
	// env := &Environment{env: make(map[string]Entry)}
	// env.registerType("char", &EntryType{Typee: Type{Alignment: 1, Size: 1}})
	// env.registerType("i8", &EntryType{Typee: Type{Alignment: 2, Size: 1}})
	// env.registerType("i16", &EntryType{Typee: Type{Alignment: 2, Size: 2}})
	// env.registerType("i32", &EntryType{Typee: Type{Alignment: 4, Size: 4}})
	// env.registerType("i64", &EntryType{Typee: Type{Alignment: 8, Size: 8}})
	// some_var := StmtLet{Name: "some_var", TypeName: "i64", Init: &ExprInt{Value: 1}}
	// other_var := StmtLet{Name: "other_var", TypeName: "i64", Init: &ExprInt{Value: 1}}
	// other_one := StmtLet{Name: "other_one", TypeName: "i64", Init: &ExprInt{Value: 4}}
	// // left := ExprInt{Value: 1}
	// right := ExprInt{Value: 2}
	// ident := ExprIdent{Name: "some_var"}
	// ident2 := ExprIdent{Name: "other_var"}
	// ident3 := ExprIdent{Name: "other_one"}
	// bi := ExprBin{Left: &ident2, Right: &right, Op: '+'}
	// bi2 := ExprBin{Left: &ident, Right: &ident3, Op: '*'}
	// resolver(env, &some_var)
	// resolver(env, &other_var)
	// resolver(env, &other_one)
	// resolver(env, &bi2)
	// resolver(env, &bi)
	//
	// nodes := []Node{&some_var, &other_one, &other_var, &bi, &bi2}
	// codegen(env, nodes)
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s path-to-file\n", os.Args[0])
		return
	}
	filePath := os.Args[1]
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("Provided file '%s' doesn't exist.", filePath)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		println("Unable to open file")
		return
	}
	src, err := io.ReadAll(file)
	file.Close()
	bag := error.New()
	lex := lexer.New(bag)
	tokens := lex.GetTokens([]byte(src))
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	parser := parser.New(tokens, bag)
	tree := parser.Parse()
	resolver.Resolve(tree, bag)
	if bag.GotErrors() {
		bag.PrintErrors()
		os.Exit(1)
	}

}
