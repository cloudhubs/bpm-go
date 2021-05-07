# RAD GO

REST API Discovery for Go applications.

## Build and run

```
$ goimports -l -w ./
$ go run main.go
```

Request:
```
POST 127.0.0.1:8085/parse
{
  "filePath": "/Users/das/Downloads/main.go"
}
```

Response:
```
{
  "filePath": "/Users/das/Downloads/main.go"
}
```

## Authors

- [Dipta Das](https://github.com/diptadas)
- [Maruf Tuhin](https://github.com/the-redback)
