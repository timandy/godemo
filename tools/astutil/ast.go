package astutil

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

func GetFieldType(fieldList []*ast.Field, name string) string {
	for _, field := range fieldList {
		if len(field.Names) == 0 {
			continue
		}
		ident := field.Names[0]
		if ident == nil {
			continue
		}
		if ident.Name != name {
			continue
		}
		typeIdent := field.Type.(*ast.Ident)
		if typeIdent == nil {
			continue
		}
		return typeIdent.Name
	}
	return ""
}

func CreateField(name string, typ string) *ast.Field {
	return &ast.Field{Names: []*ast.Ident{ast.NewIdent(name)}, Type: ast.NewIdent(typ)}
}

func IndexAssignTimerNil(bodyList []ast.Stmt) (x ast.Expr, index int) {
	for idx, stmt := range bodyList {
		// check is AssignStmt
		as, isAs := stmt.(*ast.AssignStmt)
		if !isAs {
			continue
		}
		// check is assign operator
		if as.Tok != token.ASSIGN {
			continue
		}
		// check the length of lhs is 1
		lhs := as.Lhs
		if len(lhs) != 1 {
			continue
		}
		// check lhs[0] is SelectorExpr
		se, isSe := lhs[0].(*ast.SelectorExpr)
		if !isSe {
			continue
		}
		// check selector is ?.timer
		sel := se.Sel
		if sel == nil || sel.Name != "timer" {
			continue
		}
		// check the length of rhs is 1
		rhs := as.Rhs
		if len(rhs) != 1 {
			continue
		}
		// check rhs[0] is Ident
		ident, isIdent := rhs[0].(*ast.Ident)
		if !isIdent {
			continue
		}
		// check the AssignStmt is ?.timer = nil
		if ident.Name != "nil" {
			continue
		}
		return se.X, idx
	}
	return nil, -1
}

func CreateAssignNilStmt(x ast.Expr, name string) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.SelectorExpr{X: x, Sel: &ast.Ident{Name: name}}},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{&ast.Ident{Name: "nil"}},
	}
}

func SaveAs(sourcePath string, destDir string, fset *token.FileSet, af *ast.File) string {
	shortName := filepath.Base(sourcePath)
	destPath := filepath.Join(destDir, shortName)
	// create dest file
	destFile, err := os.Create(destPath)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()
	// write code to dest file
	err = printer.Fprint(destFile, fset, af)
	if err != nil {
		panic(err)
	}
	return destPath
}
