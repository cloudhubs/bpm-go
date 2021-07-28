package lib

import (
	"fmt"
)

func RunProjectAnalysis(request ParseRequest) ([]FunctionNode, error) {
	if request.Path == "" || request.ProjectKey == "" {
		return nil, fmt.Errorf("missing path or project key")
	}

	fnNodes, err := GetFunctionNodes(request)
	if err != nil {
		return nil, err
	}

	sonarRes, err := RunSonarAnalysis(request)
	if err != nil {
		return nil, err
	}

	// attach issues to nodes
	for i, node := range fnNodes {
		for _, issue := range sonarRes.Issues {
			if node.FilePath == issue.FilePath && node.Name == issue.Function {
				fnNodes[i].Issues = append(fnNodes[i].Issues, issue)
			}
		}
	}

	return fnNodes, nil
}
