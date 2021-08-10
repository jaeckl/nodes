package main

import (
    "github.com/jaeckl/nodes/core"
    "github.com/jaeckl/nodes/pack"
)


func main() {
    rt := core.NewRuntime()
    pack.Load(rt,"core.npk")

    rt.New("core/math/add.nd:1","core/math/add.nd","")
    rt.New("core/lang/value.nd:1","core/lang/value.nd","int 5")
    rt.New("core/lang/value.nd:2","core/lang/value.nd","int 0")

    rt.Connect("core/lang/value.nd:1","out","core/math/add.nd:1","left")
    rt.Connect("core/lang/value.nd:1","out","core/math/add.nd:1","right")
    rt.Connect("core/math/add.nd:1","out","core/lang/value.nd:2","in")
    rt.ConnectMsg("core/math/add.nd:1","core/lang/value.nd:2")

    rt.ReceiveMessage("core/math/add.nd:1","Pulse")
}




