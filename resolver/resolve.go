package resolver

import (
	"fmt"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/types"
)

type Table struct {
	Symbols *scope.Scope
}

const POINTER_SIZE = 8
const POINTER_ALIGNMENT = 8

var table *Table
var handler *error.DiagnosticBag
var cachedPtrTypes map[string]*types.Type = make(map[string]*types.Type)

func InitTable() *Table {
	t := Table{Symbols: scope.NewScope(nil)}
	t.Symbols.Define("i8", scope.NewTypeObj(types.NewType("i8", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("i16", scope.NewTypeObj(types.NewType("i16", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("i32", scope.NewTypeObj(types.NewType("i32", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("i64", scope.NewTypeObj(types.NewType("i64", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("bool", scope.NewTypeObj(types.NewType("bool", types.TYPE_BOOL, 1, 1)))
	t.Symbols.Define("void", scope.NewTypeObj(types.NewType("void", types.TYPE_VOID, 0, 0)))
	t.Symbols.Define("string", scope.NewTypeObj(types.NewType("string", types.TYPE_STRING, 0, 0)))
	return &t
}
func Resolve(program []ast.Decl, bag *error.DiagnosticBag) (*Table, []DeclNode) {
	println("----RESOLVER----")
	table = InitTable()
	handler = bag
	var decls []DeclNode
	for _, decl := range program {
		decls = append(decls, resolveDecl(decl))
	}
	if !table.Symbols.Lookup("main") {
		handler.ReportError(error.Position{}, "'main' function couldn't be found")
	}
	return table, decls
}
func isTypeExist(typee ast.TypeSpec) (*types.Type, bool) {
	switch t := typee.(type) {
	case *ast.TypeName:
		{
			if table.Symbols.Lookup(t.Name) {
				obj := table.Symbols.GetObj(t.Name)
				if obj.Kind != scope.TYPE {
					handler.ReportError(typee.GetPos(), "Type '%s' must be a type not variable name or function name", t.Name)
					return nil, false
				}
				return obj.Type, true
			}
			handler.ReportError(typee.GetPos(), "Type '%s' doesn't exist", t.Name)
		}
	case *ast.TypePtr:
		{
			val, ok := isTypeExist(t.Base)
			if ok {
				typeName := "*" + val.TypeName
				if cachedPtrTypes[typeName] != nil {
					return cachedPtrTypes[typeName], true
				}
				ptr := types.NewType(typeName, types.TYPE_PTR, POINTER_SIZE, POINTER_ALIGNMENT)
				ptr.Base = val
				cachedPtrTypes[typeName] = ptr
				return ptr, true
			}
		}
	}
	return nil, false
}
func resolveDecl(decl ast.Decl) DeclNode {
	switch node := decl.(type) {
	case *ast.DeclFunction:
		{
			if table.Symbols.LookupOnce(node.Name) {
				handler.ReportError(node.Pos, "Can't redeclare function '%s' more than once", node.Name)
				return nil
			}
			retType, ok := isTypeExist(node.RetType)
			if !ok {
				return nil
			}
			fnScope := scope.NewScope(nil)
			table.Symbols.Define(node.Name, scope.NewObj(scope.FN, retType))
			for _, param := range node.Parameters {
				if !fnScope.LookupOnce(param.Name) {
					typ, ok := isTypeExist(param.Type)
					if !ok {
						return nil
					}
					fnScope.Define(param.Name, scope.NewObj(scope.PARAM, typ))
				} else {
					handler.ReportError(node.Pos, "Can't redeclare '%s' parameter more than once", param.Name)
					return nil
				}
			}
			table.Symbols.GetObj(node.Name).Scope = fnScope
			resolvedBody := resolveStmt(node.Body, fnScope)
			return &DeclFunction{Scope: fnScope, Name: node.Name, Body: resolvedBody, ReturnType: retType}
		}
	case *ast.DeclStruct:
		{
			if table.Symbols.LookupOnce(node.Name) {
				handler.ReportError(node.Pos, "Can't redeclare struct '%s' more than once", node.Name)
				return nil
			}
			structScope := scope.NewScope(nil)
			obj := scope.NewObj(scope.TYPE, types.NewType(node.Name, types.TYPE_STRUCT, 0, 0))
			obj.Scope = structScope
			table.Symbols.Define(node.Name, obj)
			fields := make([]Field, 0, 4)
			for _, field := range node.Fields {
				if structScope.LookupOnce(field.Name) {
					handler.ReportError(field.Pos, "Can't redeclare '%s' field more than once in struct '%s'", field.Name, node.Name)
					return nil
				}
				typ, ok := isTypeExist(field.Type)
				if !ok {
					return nil
				}
				obj := scope.NewObj(scope.FIELD, typ)
				if typ.Kind == types.TYPE_STRUCT {
					obj.Scope = table.Symbols.GetObj(typ.TypeName).Scope
				}
				structScope.Define(field.Name, obj)
				fields = append(fields, Field{Name: field.Name, Type: typ})
			}
			return &DeclStruct{Name: node.Name, Fields: fields, Pos: node.Pos, Scope: structScope}
		}
	}
	return nil
}

func resolveStmt(stmt ast.Stmt, currScope *scope.Scope) StmtNode {
	pos := stmt.GetPos()
	switch node := stmt.(type) {
	case *ast.StmtLet:
		{
			if !currScope.LookupOnce(node.Name) {
				typ, ok := isTypeExist(node.Type)
				if !ok {
					return nil
				}
				currScope.Define(node.Name, scope.NewObj(scope.VAR, typ))
				var resolvedExpr ExprNode
				if node.Init != nil {
					resolvedExpr = resolveExpr(node.Init, currScope, nil)
				}
				return &StmtLet{Name: node.Name, Init: resolvedExpr, Scope: currScope, Type: typ, Pos: node.Pos}
			}
			handler.ReportError(pos, "Can't redeclare '%s' variable more than once in same block", node.Name)
		}
	case *ast.StmtReturn:
		{
			if node.Result != nil {
				resolvedExpr := resolveExpr(node.Result, currScope, nil)
				return &StmtReturn{Result: resolvedExpr}
			}
		}
	case *ast.StmtExpr:
		{
			expr := resolveExpr(node.Expr, currScope, nil)
			return &StmtExpr{Expr: expr, Scope: currScope}
		}
	case *ast.StmtBlock:
		{
			s := scope.NewScope(currScope)
			var resolvedStmts []StmtNode
			for _, stmt := range node.Block {
				resolvedStmts = append(resolvedStmts, resolveStmt(stmt, s))
			}
			return &StmtBlock{Scope: s, Body: resolvedStmts}
		}
	}
	return nil
}

func resolveExpr(expr ast.Expr, currScope *scope.Scope, typeScope *scope.Scope) ExprNode {
	pos := expr.GetPos()
	switch node := expr.(type) {
	case *ast.ExprBinary:
		{
			left := resolveExpr(node.Left, currScope, nil)
			right := resolveExpr(node.Right, currScope, nil)
			return &ExprBinary{Left: left, Right: right, Op: KindToBinary[node.Op]}
		}
	case *ast.ExprCompound:
		{
			typ, ok := isTypeExist(node.Type)
			if !ok {
				return nil
			}
			// Type must not be a pointer or primitive  e.g '*Vector{}','i32'
			if typ.Kind != types.TYPE_STRUCT /* || union*/ {
				handler.ReportError(node.Pos, "Type must be a struct or union in order to compose")
				return nil
			}
			structScope := table.Symbols.GetObj(typ.TypeName).Scope
			fieldsName := structScope.QueryByKind(scope.FIELD)
			resolvedFieldsName := map[string]bool{}
			resolvedFields := make([]ExprCompoundField, 0, 4)
			for _, field := range node.Fields {
				if !structScope.LookupOnce(field.Name) {
					handler.ReportError(field.Pos, "'%s' doesn't have '%s' field", typ.TypeName, field.Name)
					continue
				}
				resolvedFieldsName[field.Name] = true
				resolvedExpr := resolveExpr(field.Init, currScope, nil)
				resolvedFields = append(resolvedFields, ExprCompoundField{Name: field.Name, Expr: resolvedExpr})
			}
			for _, fieldName := range fieldsName {
				_, ok := resolvedFieldsName[fieldName]
				if !ok {
					handler.ReportError(node.Pos, "Field '%s' must be initialized in '%s' struct Compound", fieldName, typ.TypeName)
				}
			}
			return &ExprCompound{Type: typ, Fields: resolvedFields, Pos: node.Pos}
		}
	case *ast.ExprCall:
		{
			if !table.Symbols.LookupOnce(node.Name) {
				handler.ReportError(node.Pos, "Function '%s' not found", node.Name)
				return nil
			}
			fnObj := table.Symbols.GetObj(node.Name)
			params := fnObj.Scope.QueryByKind(scope.PARAM)
			args := make([]*ExprArg, 0)
			for _, arg := range node.Args {
				resolved := resolveExpr(arg, currScope, nil)
				args = append(args, &ExprArg{Expr: resolved})
			}
			paramLen := len(params)
			argsLen := len(args)
			if paramLen != argsLen {
				handler.ReportError(node.Pos, "Function '%s' expected '%d' arguments but got '%d' arguments", node.Name, paramLen, argsLen)
				return nil
			}
			return &ExprCall{Name: node.Name, Args: args, Pos: node.Pos}
		}
	case *ast.ExprAssign:
		{
			left := resolveExpr(node.Left, currScope, nil)
			right := resolveExpr(node.Right, currScope, nil)
			return &ExprAssign{Right: right, Left: left}
		}
	case *ast.ExprInt:
		{
			return &ExprInt{Value: node.Value}
		}
	case *ast.ExprBoolean:
		{
			return &ExprBool{Value: node.Value}
		}
	case *ast.ExprString:
		{
			return &ExprString{Value: node.Value}
		}

	case *ast.ExprUnary:
		{
			resolved := resolveExpr(node.Right, currScope, nil)
			if resolved != nil {
				return &ExprUnary{Type: resolved.GetType(), Right: resolved, Op: KindToUnary[node.Op]}
			}
		}
	case *ast.ExprField:
		{
			left := resolveExpr(node.Expr, currScope, nil)
			if left != nil {
				typ := left.GetType()
				if typ.Kind == types.TYPE_PTR {
					typ = typ.Base
				}
				typeName := typ.TypeName
				if typ.Kind != types.TYPE_STRUCT {
					handler.ReportError(left.GetPos(), "Primitive type '%s' doesn't have fields", typeName)
					return nil
				}
				structScope := table.Symbols.GetObj(typeName).Scope
				if structScope.LookupOnce(node.Name) {
					return &ExprField{Type: structScope.GetObj(node.Name).Type, Name: node.Name}
				} else {
					handler.ReportError(left.GetPos(), "'%s' doesn't have '%s' field", typeName, node.Name)
				}
			}
			return left
		}

	case *ast.ExprIdent:
		{
			if !currScope.Lookup(node.Name) {
				handler.ReportError(pos, "Variable '%s' not found", node.Name)
				return nil
			}
			return &ExprIdentifier{Name: node.Name, Type: currScope.GetObj(node.Name).Type}
		}
	default:
		{
			panic(fmt.Sprintf("Unhandled node '%T' or Unreachable", node))
		}
	}
	return nil
}
