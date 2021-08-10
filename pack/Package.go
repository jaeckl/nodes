package pack
import "log"
import "archive/zip"
import "plugin"
import "path/filepath"
import "os"
import "io"
import "github.com/jaeckl/nodes/core"
import "github.com/jaeckl/nodes/api"

const (
    ext = ".nd"
)

func Load(rt *core.Runtime,path string) {
    reader,err := zip.OpenReader(path)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    for _,file := range reader.File {
        if !file.FileHeader.FileInfo().IsDir() && filepath.Ext(file.Name) == ext {
            rc, err := file.Open()
            defer rc.Close()

            /** Move it out of the zip to load*/
            path := filepath.Join("/tmp",filepath.Base(file.Name))
            temp, err := os.Create(path)
            io.Copy(temp,rc)
            defer temp.Close()

            p,err := plugin.Open(path)
            if err != nil {
                panic(err)
            }

            n, err := p.Lookup("New")
            if err != nil {
                panic(err)
            }
            rt.RegisterObject(file.Name,n.(func()api.NodeObject))
        }
    }
}
