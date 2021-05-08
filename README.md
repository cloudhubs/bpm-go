# RAD GO

REST API Discovery for Go applications.

## Build and run

```
$ goimports -l -w ./
$ go run main.go
```

Request:
```
POST 127.0.0.1:8085/importContext
{
    "path": "/Users/das/Baylor/RA/rad-go"
}
```

Response:
```JSON
{
    "rootPath": "/Users/das/Baylor/RA/rad-go",
    "imports": [
        {
            "path": "/Users/das/Baylor/RA/rad-go/app/app.go",
            "packages": [
                "log",
                "net/http",
                "rad-go/app/handler",
                "github.com/gorilla/mux"
            ]
        },
        {
            "path": "/Users/das/Baylor/RA/rad-go/app/handler/common.go",
            "packages": [
                "encoding/json",
                "net/http"
            ]
        },
        {
            "path": "/Users/das/Baylor/RA/rad-go/app/handler/parser.go",
            "packages": [
                "encoding/json",
                "go/parser",
                "go/token",
                "net/http",
                "os",
                "path/filepath",
                "rad-go/app/model"
            ]
        },
        {
            "path": "/Users/das/Baylor/RA/rad-go/app/model/model.go",
            "packages": null
        },
        {
            "path": "/Users/das/Baylor/RA/rad-go/main.go",
            "packages": [
                "rad-go/app"
            ]
        }
    ]
}
```

## Authors

- [Dipta Das](https://github.com/diptadas)
- [Maruf Tuhin](https://github.com/the-redback)
