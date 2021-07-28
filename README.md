# BPM GO

Technical Debt analysis for Go applications.

- Business Process Modeling (BPM)
- SonarQube analysis
- Log processing

## Build and Run

```
$ go mod tidy && goimports -l -w ./
$ go run main.go
```

## API Usage

### Project Analysis

```shell
POST 127.0.0.1:8085/analysis
{
    "path": "/Users/das/Baylor/RA/ccx-notification-service",
    "projectKey": "ccx"
}
```

### Function Calls

```shell
POST 127.0.0.1:8085/bpm
{
    "path": "/Users/das/Baylor/RA/ccx-notification-service"
}
```

### Sonar Analysis

```shell
POST 127.0.0.1:8085/sonar
{
    "path": "/Users/das/Baylor/RA/ccx-notification-service",
    "projectKey": "ccx"
}
```
