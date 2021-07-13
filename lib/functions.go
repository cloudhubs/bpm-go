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
	parseZeroLog(node)
	return v

	//switch node := node.(type) {
	//case *ast.CallExpr:
	//	switch nodeFn := node.Fun.(type) {
	//	case *ast.SelectorExpr:
	//		if pkgID, ok := nodeFn.X.(*ast.Ident); ok {
	//			if pkgID.String() == "log" {
	//				fmt.Println(pkgID)
	//				if node.Args != nil {
	//					if basicLit, ok := node.Args[0].(*ast.BasicLit); ok {
	//						msg := basicLit.Value
	//						if msg != "" {
	//							fmt.Println("log message", basicLit.Value)
	//						}
	//					}
	//				}
	//			}
	//		}
	//
	//		/*if basicLit, ok := nodeFn.X.(*ast.CallExpr).Args[0].(*ast.BasicLit); ok {
	//			fmt.Println("log message", basicLit.Value)
	//		}*/
	//
	//		v.fnCallExpr = append(v.fnCallExpr, nodeFn.Sel.Name)
	//	}
	//}

	return v
}

func parseZeroLogRec(node interface{}) []string {
	if n1, ok := node.(*ast.CallExpr); ok {
		argsVal := ""
		for _, x := range n1.Args {
			if xv, ok := x.(*ast.BasicLit); ok {
				argsVal += " " + xv.Value
			} else if xv, ok := x.(*ast.Ident); ok {
				argsVal += " " + xv.String()
			}
		}
		return append(parseZeroLogRec(n1.Fun), argsVal)
	} else if n2, ok := node.(*ast.SelectorExpr); ok {
		return append(parseZeroLogRec(n2.X), n2.Sel.String())
	} else if n3, ok := node.(*ast.Ident); ok {
		return []string{n3.String()}
	} else {
		return nil
	}
}

func parseZeroLog(node ast.Node) string {
	//n1, ok := node.(*ast.CallExpr)
	//if !ok {
	//	return ""
	//}
	//n2, ok := n1.Fun.(*ast.SelectorExpr)
	//if !ok {
	//	return ""
	//}
	//n3, ok := n2.X.(*ast.CallExpr)
	//if !ok {
	//	return ""
	//}
	//n4, ok := n3.Fun.(*ast.SelectorExpr)
	//if !ok {
	//	return ""
	//}
	//n5, ok := n4.X.(*ast.CallExpr)
	//if !ok {
	//	return ""
	//}
	//n6, ok := n5.Fun.(*ast.SelectorExpr)
	//if !ok {
	//	return ""
	//}
	//n7, ok := n6.X.(*ast.Ident)
	//if !ok {
	//	return ""
	//}
	//
	//n7Val := n7.String()
	//if n7Val != "log" {
	//	return ""
	//}
	//n6SelVal := n6.Sel.String()
	//if n6SelVal != "Info" {
	//	return ""
	//}
	//n4SelVal := n4.Sel.String()
	//if n4SelVal != "Bool" {
	//	return ""
	//}
	//n2SelVal := n2.Sel.String()
	//if n2SelVal != "Msg" {
	//	return ""
	//}
	//
	//n3Args := n3.Args
	//n3ArgsVal := ""
	//for _, x := range n3Args {
	//	if xv, ok := x.(*ast.BasicLit); ok {
	//		n3ArgsVal += " " + xv.Value
	//	} else if xv, ok := x.(*ast.Ident); ok {
	//		n3ArgsVal += " " + xv.String()
	//	}
	//}
	//
	//n1Args := n1.Args
	//n1ArgsVal := ""
	//for _, x := range n1Args {
	//	if xv, ok := x.(*ast.BasicLit); ok {
	//		n1ArgsVal += " " + xv.Value
	//	} else if xv, ok := x.(*ast.Ident); ok {
	//		n1ArgsVal += " " + xv.String()
	//	}
	//}
	//
	//fmt.Println(n7Val, n6SelVal, n4SelVal, n3ArgsVal, n2SelVal, n1ArgsVal)

	stmt := parseZeroLogRec(node)
	if len(stmt) == 7 && stmt[0] == "log" && stmt[5] == "Msg" {
		fmt.Println(stmt)
	}

	return ""
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
