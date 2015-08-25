package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const file = "../sampleServer/main.go"

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the imports from the file's AST.
	// for _, s := range f.Imports {
	// 	fmt.Println(s.Path.Value)
	// }
	// for _, d := range f.Decls {
	// 	fmt.Println(d)
	// }

	// ast.Print(fset, f.Decls[1])

	ast.Walk(new(FuncVisitor), f.Decls[1])
}

type FuncVisitor struct {
}

func (v *FuncVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.FuncDecl:
		t.Name = ast.NewIdent(strings.Title(t.Name.Name))
		fmt.Println(t.Name.String())
		fmt.Printf("%+v\n", t.Type)
	case *ast.CompositeLit:
		// t.Name = ast.NewIdent(strings.Title(t.Name.Name))
		// fmt.Println(t.Name.String())
		// fmt.Printf("CompositeLit: %+v\n", t)
		fmt.Printf("CompositeLit Type: %+v\n", t.Type)
		// fmt.Printf("CompositeLit Elts: %+v\n", t.Elts[0])
	case *ast.AssignStmt:
		// fmt.Printf("AssignStmt: %+v\n", t)
		// fmt.Printf("AssignStmt RHS: %+v\n", t.Rhs[0])
		switch r := t.Rhs[0].(type) {
		case *ast.CompositeLit:
			// fmt.Printf("rhs composite lit: %+v\n", r.Type)
			switch rt := r.Type.(type) {
			case *ast.SelectorExpr:
				// if rt.X == "snowflake" && rt.Sel == "Resource" {
				fmt.Printf("X Sel: %+v %+v\n", rt.X, rt.Sel)
				// }
			}
		}
	}

	return v
}
