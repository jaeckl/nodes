package pack
import "log"
import "archive/zip"
import "fmt"
type Package struct {
}

func Load(path string) {
    reader,err := zip.OpenReader(path)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    for _,file := range reader.File {
        fmt.Println(file.Name)
    }
}
