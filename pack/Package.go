package pack
import "log"
import "archive/zip"
import "fmt"
import "plugin"
import "path/filepath"
import "os"
import "io"

const (
    ext = ".nd"
)
type Package struct {
}

func Load(path string) {
    reader,err := zip.OpenReader(path)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    for _,file := range reader.File {
        if !file.FileHeader.FileInfo().IsDir() && filepath.Ext(file.Name) == ext {
            rc, err := file.Open()



            temp, err := os.Create("/tmp/tmp.nd")
            io.Copy(temp,rc)

            _,err = plugin.Open("/tmp/tmp.nd")
            fmt.Println(err)
            if err != nil {
                panic(err)
            }
        }
    }
}
