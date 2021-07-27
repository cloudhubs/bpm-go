package lib

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const SonarUrl = "http://localhost:9000"
const SonarUser = "admin"
const SonarPass = "admin"
const ProjectKey = "bu-project"

type SonarResult struct {
	Total  int
	Issues []SonarIssue
}

type SonarIssue struct {
	Component string
	Line      int
	Severity  string
	Rule      string
	Type      string
	Message   string
	Effort    string
	Debt      string
}

func init() {
	// TODO: run Sonar server separately?
	go startSonar()
	for ; isSonarUp() == false; {
		time.Sleep(10 * time.Second)
	}
}

// docker run --rm -p 9000:9000 sonarqube:8.2-community
func startSonar() {
	cmd := exec.Command("docker", "run", "--rm",
		"-p=9000:9000",
		"sonarqube:8.2-community",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println("run sonar error:", err)
	}
}

func isSonarUp() bool {
	healthApi := SonarUrl + "/api/system/health"
	resp, err := resty.New().SetBasicAuth(SonarUser, SonarPass).R().Get(healthApi)
	if err != nil {
		log.Println("Sonar health API error:", err)
		return false
	}
	log.Println("Sonar health API response:", resp.String())
	return strings.Contains(resp.String(), "GREEN")
}

func deleteProject() error {
	deleteApi := fmt.Sprintf("%s/api/projects/delete?project=%s", SonarUrl, ProjectKey)
	resp, err := resty.New().SetBasicAuth(SonarUser, SonarPass).R().Post(deleteApi)
	if err != nil {
		return fmt.Errorf("sonar delete project API error, reason: %v", err)
	}
	log.Println("Sonar delete project API response:", resp.String())
	return nil
}

func runSonarScanner(sourcePath string) error {
	// delete the project before running a new scan
	err := deleteProject()
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "run", "--rm",
		fmt.Sprintf("-v=%s:/usr/src", sourcePath),
		"--network=host",
		"sonarsource/sonar-scanner-cli",
		"-D", fmt.Sprintf("sonar.projectKey=%s", ProjectKey),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("sonar run scanner error, reason: %v", err)
	}

	return nil
}

func getGolangProjectIssues() (SonarResult, error) {
	// TODO: loop through pages
	issuesApi := fmt.Sprintf("%s/api/issues/search?lcomponentKeys=%s&languages=go&ps=500&p=1", SonarUrl, ProjectKey)
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

func RunSonarAnalysis(sourcePath string) (SonarResult, error) {
	err := runSonarScanner(sourcePath)
	if err != nil {
		return SonarResult{}, err
	}
	return getGolangProjectIssues()
}
