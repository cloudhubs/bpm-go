# BPM GO

Business Process Modeling for Go applications.

## Build and run

```
$ go mod tidy && goimports -l -w ./
$ go run main.go
```

Request:
```
POST 127.0.0.1:8085/imports
{
    "path": "/Users/das/Baylor/RA/bpm-go"
}
```

Response:
```JSON
{
  "rootPath": "/Users/das/Baylor/RA/bpm-go",
  "imports": [
    {
      "path": "/Users/das/Baylor/RA/bpm-go/api/common.go",
      "packages": [
        "encoding/json",
        "log",
        "net/http"
      ]
    },
    {
      "path": "/Users/das/Baylor/RA/bpm-go/api/handler.go",
      "packages": [
        "bpm-go/lib",
        "encoding/json",
        "net/http"
      ]
    },
    {
      "path": "/Users/das/Baylor/RA/bpm-go/api/server.go",
      "packages": [
        "fmt",
        "log",
        "net/http",
        "github.com/gorilla/mux"
      ]
    },
    {
      "path": "/Users/das/Baylor/RA/bpm-go/lib/helper.go",
      "packages": [
        "os",
        "path/filepath"
      ]
    },
    {
      "path": "/Users/das/Baylor/RA/bpm-go/lib/model.go",
      "packages": null
    },
    {
      "path": "/Users/das/Baylor/RA/bpm-go/lib/process.go",
      "packages": [
        "go/parser",
        "go/token"
      ]
    },
    {
      "path": "/Users/das/Baylor/RA/bpm-go/main.go",
      "packages": [
        "bpm-go/api"
      ]
    }
  ]
}
```

