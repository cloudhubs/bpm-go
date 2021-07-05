# BPM GO

Business Process Modeling for Go applications.

## Build and Run

```
$ go mod tidy && goimports -l -w ./
$ go run main.go
```

## API Usage

### Imports

```
POST 127.0.0.1:8085/imports
{
    "path": "/Users/das/Baylor/RA/bpm-go"
}
```

### Function Calls

```
POST 127.0.0.1:8085/functions
{
    "path": "/Users/das/Baylor/RA/bpm-go"
}
```