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
		fnCalls, log := funcCallsInFunc(fnNode)

		fnNodes[i].ChildNodeIDs = getChildNodeIDs(fnCalls, fnNodes)
		fnNodes[i].Logs = log
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

func funcCallsInFunc(function FunctionNode) ([]string, []*Log) {
	visitor := &FnCallVisitor{}
	visitor1 := &FnCallExpr{}
	ast.Walk(visitor1, function.funcDecl.Body)


	ast.Walk(visitor, function.funcDecl.Body)
	return visitor.fnCalls, visitor1.fnCallExpr



}

type FnCallVisitor struct {
	fnCalls []string
}

type FnCallExpr struct {
	fnCallExpr []*Log
}

func (v *FnCallExpr) Visit(node ast.Node) (w ast.Visitor) {
	log := parseZeroLog(node)
	if log != nil && len(log.LogMsg)>0{
		v.fnCallExpr = append(v.fnCallExpr, log)
		//fmt.Println("log type = ",log.Type,"logmsg = ",log.LogMsg)
	}
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

func parseZeroLog(node ast.Node) *Log {
	stmt := parseZeroLogRec(node)

	// zero log format
	if len(stmt) == 7 && stmt[0] == "log" && stmt[5] == "Msg" {
		//fmt.Println("zero log:", stmt)
		return &Log{
			Type:   stmt[1], // Info, Err, etc
			LogMsg: stmt[6],
		}
	}

	// default go log format
	if len(stmt) == 3 && stmt[0] == "log" {

		return &Log{
			Type:   stmt[1], // Print, Fatal
			LogMsg: stmt[2],
		}
	}

	return nil
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
