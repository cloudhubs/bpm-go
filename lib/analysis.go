package lib

import (
	"fmt"
)

func RunProjectAnalysis(request ParseRequest) (SonarResult, error) {
	if request.Path == "" || request.ProjectKey == "" {
		return SonarResult{}, fmt.Errorf("missing path or project key")
	}

	return SonarResult{}, nil
}
