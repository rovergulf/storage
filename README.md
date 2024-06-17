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

    fs, err := storage.NewStorage(
        storage.WithBackends(storage.LocalBackends),
        storage.WithPathPrefix(storagePath),
    )
    if err != nil {
        // log.Fatal(err)	
    }

    fileData := []byte(`{"example": "hello world"}`)
    fileName := "example.json"
	
    // check if file key already exists 
    if exists, _ := fs.Exists(ctx, fileName); exists {
        // log.Info("already exists")
        return
    }
	
	// upload data
    uploadErr := fs.Put(ctx, fileName, fileData)
    if uploadErr != nil {
        // log.Errorf("Unable to upload to '%s': %s", fileName, err)
        // handle uploadErr
    }
}
```
