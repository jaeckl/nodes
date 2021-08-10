//go:generate go build -buildmode=plugin -o node/lang/value.nd node/lang/Value.go
package main
import "github.com/jaeckl/nodes/api"
import "strconv"
import "fmt"
import "strings"
type NumberObject struct {
    dataInt int64
    dataFloat float64
    tp string
}

func New() api.NodeObject {
    object := &NumberObject{}
    return object
}

func (ad *NumberObject) Init(ctx api.RuntimeContext,args string) bool {
    argArray := strings.Split(args," ")
    if len(argArray) == 1 {
        ad.tp = argArray[0]
    } else if len(argArray) == 2 {
        ad.tp = argArray[0]
        if argArray[0] == "int" {
            if s, err := strconv.ParseInt(argArray[1], 10, 32); err == nil {
                ad.dataInt = s
            }
        }else if argArray[0] == "float" {
            if s, err := strconv.ParseFloat(argArray[1], 32); err == nil {
                ad.dataFloat = s
            }
        } else {
            return false
        }
    } else {
        return false
    }

    ctx.NewInputNode("in",ad.tp)
    ctx.NewOutputNode("out",ad.tp)

    if ad.tp == "int" {
        ctx.Outputs("out").SetValue(ad.dataInt)
    } else {
        ctx.Outputs("out").SetValue(ad.dataFloat)
    }

    return true
}

func (ad *NumberObject) ReceiveMessage(ctx api.RuntimeContext,msg string) {
    if ad.tp == "int" {
        fmt.Println(ctx.Inputs("in").GetValue().(int64))
    }
    if ad.tp == "float" {
        fmt.Println(ctx.Inputs("in").GetValue().(float64))
    }
}
