![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/rovergulf/storage)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rovergulf/storage/test?branch=main)

# storage

This repository contains source code for file storage, 
supporting local file storage and cloud remote ones.

### Install module
```shell
go get github.com/rovergulf/storage
```

### Example
```go
package main

import (
	"context"
	"log"
	
	"github.com/rovergulf/storage"
)

func main() {
    ctx := context.Background()

    storagePath := "/temp/my-storage"

    fs, err := storage.NewStorage(
        storage.WithBackends(storage.LocalBackends),
        storage.WithPrefix(storagePath),
    )
    if err != nil {
        log.Fatal(err)	
    }

    fileData := []byte(`{"example": "hello world"}`)
    fileName := "example.json"
	
    // check if file key already exists 
    if exists, _ := fs.Exists(ctx, fileName); exists {
        log.Printf("'%s' already exists", fileName)
        return
    }
	
	// upload data
    uploadErr := fs.Put(ctx, fileName, fileData)
    if uploadErr != nil {
        log.Fatalf("Unable to upload to '%s': %s", fileName, err)
        // handle uploadErr
    }
}
```

### OpenTelemetry support
This package supports [OpenTelemetry](https://opentelemetry.io) tracer spans.   
Visit [opentelemetry-go examples](https://github.com/open-telemetry/opentelemetry-go/tree/main/example) to install tracer provider and provide tracer as storage `storage.WithTracer` option
```go
package main

import (
	...
	
	"go.opentelemetry.io/otel"
	
	"github.com/rovergulf/storage"
)

func main() {
	tracer := otel.Tracer("go-storage")
	
	fs, err := storage.NewStorage(
		...,
		storage.WithTracer(tracer),
	)
	if err != nil {
		// handle err	
	}

}

```
