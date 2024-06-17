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

    fileStorage := storage.NewFileStorage(storagePath)

    fileData := []byte(`{"example": "hello world"}`)
    fileName := "example.json"
	
    // check if file key already exists 
    if exists, _ := fileStorage.Exists(ctx, fileName); exists {
        return
    }
	
	// upload data
    uploadErr := fileStorage.Put(ctx, fileName, fileData)
    if uploadErr != nil {
        // handle uploadErr
    }
}
```
