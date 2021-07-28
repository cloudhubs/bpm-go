package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const SonarUrl = "http://localhost:9000"
const SonarUser = "admin"
const SonarPass = "admin"

func init() {
	if isSonarUp() == false {
		if err := startSonar(); err != nil {
			log.Fatalln(err)
		}
	}
}

// run sonar in background
// docker run -d -p 9000:9000 sonarqube:8.2-community
func startSonar() error {
	log.Println("starting sonar server")
	cmd := exec.Command("docker", "run", "-d",
		"-p=9000:9000",
		"sonarqube:8.2-community",
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("start sonar error, reason: %v", err)
	}

	log.Println("waiting for sonar server to start")
	for isSonarUp() == false {
		time.Sleep(10 * time.Second)
	}

	return nil
}

// http://localhost:9000/api/system/health
func isSonarUp() bool {
	healthApi := SonarUrl + "/api/system/health"
	resp, err := resty.New().SetBasicAuth(SonarUser, SonarPass).R().Get(healthApi)
	if err != nil {
		log.Println("sonar health API error:", err)
		return false
	}
	log.Println("sonar health API response:", resp.String())
	return strings.Contains(resp.String(), "GREEN")
}

// http://localhost:9000/api/projects/delete?project=ccx
func deleteProject(projectKey string) error {
	deleteApi := fmt.Sprintf("%s/api/projects/delete?project=%s", SonarUrl, projectKey)
	resp, err := resty.New().SetBasicAuth(SonarUser, SonarPass).R().Post(deleteApi)
	if err != nil {
		return fmt.Errorf("sonar delete project API error, reason: %v", err)
	}
	log.Println("sonar delete project API response:", resp.String())
	return nil
}

// http://localhost:9000/api/projects/search?projects=ccx
func isProjectExists(projectKey string) (bool, error) {
	searchApi := fmt.Sprintf("%s/api/projects/search?projects=%s", SonarUrl, projectKey)
	resp, err := resty.New().SetBasicAuth(SonarUser, SonarPass).R().Get(searchApi)
	if err != nil {
		return false, fmt.Errorf("sonar search project API error, reason: %v", err)
	}
	log.Println("sonar search project API response:", resp.String())
	return strings.Contains(resp.String(), projectKey), nil
}

func runSonarScanner(sourcePath, projectKey string) error {
	cmd := exec.Command("docker", "run", "--rm",
		fmt.Sprintf("-v=%s:/usr/src", sourcePath),
		"--network=host",
		"sonarsource/sonar-scanner-cli",
		"-D", fmt.Sprintf("sonar.projectKey=%s", projectKey),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sonar run scanner error, reason: %v", err)
	}

	return nil
}

// http://localhost:9000/api/issues/search?lcomponentKeys=ccx&languages=go&ps=500&p=1
func getGolangProjectIssues(projectKey string) (SonarResult, error) {
	// TODO: loop through pages
	issuesApi := fmt.Sprintf("%s/api/issues/search?lcomponentKeys=%s&languages=go&ps=500&p=1", SonarUrl, projectKey)
	resp, err := resty.New().SetBasicAuth(SonarUser, SonarPass).R().Post(issuesApi)
	if err != nil {
		return SonarResult{}, fmt.Errorf("sonar search issues API error, reason: %v", err)
	}

	var result SonarResult

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return SonarResult{}, fmt.Errorf("failed to decode API response, reason: %v", err)
	}
	defer resp.RawBody().Close()

	return result, nil
}

func RunSonarAnalysis(request ParseRequest) (SonarResult, error) {
	if request.Path == "" || request.ProjectKey == "" {
		return SonarResult{}, fmt.Errorf("missing path or project key")
	}

	// scan if project do not exists
	if ok, err := isProjectExists(request.ProjectKey); err != nil {
		return SonarResult{}, err
	} else if !ok {
		if err = runSonarScanner(request.Path, request.ProjectKey); err != nil {
			return SonarResult{}, err
		}
	}

	result, err := getGolangProjectIssues(request.ProjectKey)
	if err != nil {
		return SonarResult{}, err
	}

	for i, issue := range result.Issues {
		// resolve file path
		fileName := strings.TrimPrefix(issue.Component, request.ProjectKey+":")
		path := filepath.Join(request.Path, fileName)

		// resolve function name
		fn, err := findFunctionName(path, issue.Line)
		if err != nil {
			return SonarResult{}, err
		}

		result.Issues[i].FilePath = path
		result.Issues[i].Function = fn
	}

	return result, nil
}

func findFunctionName(path string, line int) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	curFunc := "global"

	for curLine, scanner := 0, bufio.NewScanner(file); scanner.Scan(); curLine++ {
		text := strings.TrimLeft(scanner.Text(), " ")
		if strings.HasPrefix(text, "func") {
			text = strings.TrimPrefix(text, "func")
			text = strings.TrimLeft(text, " ")
			curFunc = strings.SplitN(text, "(", 2)[0]
			curFunc = strings.Trim(curFunc, " ")
		}
		if curLine == line {
			break
		}
	}

	return curFunc, nil
}
