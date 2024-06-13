# storage

This repository contains source code for file storage, 
supporting local file storage and cloud remote ones.

### Install module
```shell
go get github.com/rovergulf/storage
```

### Example
```go
func main() {
    ctx := context.Background()

    storagePath := "/temp/my-storage"

    fileStorage := NewFileStorage(storagePath)

    fileData := []byte(`{"example": "hello world"}`)
    fileName := example.json

    err := fileStorage.Put(ctx, fileName, fileData)
	if err != nil {
		// handle err
    }
}
```
