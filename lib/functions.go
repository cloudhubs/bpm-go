package lib

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type AstFileWrapper struct {
	filePath string
	astFile  *ast.File
}

func GetFunctionNodes(request ParseRequest) ([]FunctionNode, error) {
	asts, err := getAsts(request.Path)

	if err != nil {
		return nil, err
	}

	var fnNodes []FunctionNode
	for _, astWrapper := range asts {
		fns := funcDeclarationsInAst(astWrapper)

		fnNodes = append(fnNodes, fns...)
	}

	for i, fnNode := range fnNodes {
		fmt.Println(fnNode)
		fmt.Println("------------------")
		fnCalls := funcCallsInFunc(fnNode)

		fnNodes[i].ChildNodeIDs = getChildNodeIDs(fnCalls, fnNodes)
	}

	return fnNodes, nil
}

// ast for each in go file in the path
func getAsts(path string) ([]AstFileWrapper, error) {
	// list all go file
	sources, err := WalkMatch(path, "*.go")
	//fmt.Println(sources)
	if err != nil {
		return nil, err
	}

	// ast for each file
	var asts []AstFileWrapper
	for _, source := range sources {
		astFile, err := parser.ParseFile(token.NewFileSet(), source, nil, 0)
		if err != nil {
			return nil, err
		}
		asts = append(asts, AstFileWrapper{
			filePath: source,
			astFile:  astFile,
		})
	}

	return asts, nil
}

var fnNodeID = 0

func getNextID() int {
	fnNodeID++
	return fnNodeID
}

func funcDeclarationsInAst(astWrapper AstFileWrapper) []FunctionNode {
	var fnNodes []FunctionNode

	for _, f := range astWrapper.astFile.Decls {
		if fn, ok := f.(*ast.FuncDecl); ok {
			fnNodes = append(fnNodes, FunctionNode{
				ID:       getNextID(),
				Name:     fn.Name.Name,
				Package:  astWrapper.astFile.Name.Name,
				FilePath: astWrapper.filePath,
				funcDecl: fn,
			})
		}
	}

	return fnNodes
}

func funcCallsInFunc(function FunctionNode) []string {
	visitor := &FnCallVisitor{}
	visitor1 := &FnCallExpr{}
	ast.Walk(visitor1, function.funcDecl.Body)
	for i, _ := range visitor1.fnCallExpr {
		i = i + 1
	}

	ast.Walk(visitor, function.funcDecl.Body)
	return visitor.fnCalls

	//var fnCalls []string
	//for _, stmt := range function.funcDecl.Body.List {
	//	if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
	//		if call, ok := exprStmt.X.(*ast.CallExpr); ok {
	//			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
	//				fnCalls = append(fnCalls, fun.Sel.Name)
	//			}
	//		}
	//	}
	//}
	//return fnCalls

}

type FnCallVisitor struct {
	fnCalls []string
}

type FnCallExpr struct {
	fnCallExpr []string
}

func (v *FnCallExpr) Visit(node ast.Node) (w ast.Visitor) {
	switch node := node.(type) {
	case *ast.CallExpr:

		switch nodeFn := node.Fun.(type) {
		case *ast.SelectorExpr:
			if pkgID, ok := nodeFn.X.(*ast.Ident); ok {
				if pkgID.String() == "log" {
					fmt.Println(pkgID)
					if node.Args !=nil {
						if basicLit, ok := node.Args[0].(*ast.BasicLit); ok {
							msg := basicLit.Value
							if msg != "" {
								fmt.Println("log message", basicLit.Value)
							}
						}
					}
				}
			}




			/*if basicLit, ok := nodeFn.X.(*ast.CallExpr).Args[0].(*ast.BasicLit); ok {
				fmt.Println("log message", basicLit.Value)
			}*/

			v.fnCallExpr = append(v.fnCallExpr, nodeFn.Sel.Name)
		}
	}

	return v
}

func (v *FnCallVisitor) Visit(node ast.Node) (w ast.Visitor) {
	// TODO: https://stackoverflow.com/questions/55377694/how-to-find-full-package-import-from-callexpr
	switch node := node.(type) {
	case *ast.CallExpr:
		switch node := node.Fun.(type) {
		case *ast.SelectorExpr:
			v.fnCalls = append(v.fnCalls, node.Sel.Name)
		case *ast.Ident:
			v.fnCalls = append(v.fnCalls, node.Name)
		}
	}

	return v
}

func getChildNodeIDs(fnCalls []string, fnNodes []FunctionNode) []int {
	var childNodeIDs []int

	for _, fnCall := range fnCalls {
		for _, fnNode := range fnNodes {
			// TODO: match package
			if fnCall == fnNode.Name {
				childNodeIDs = append(childNodeIDs, fnNode.ID)
				break
			}
		}
	}

	return childNodeIDs
}
