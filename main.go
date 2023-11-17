package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"mygf/internal/cmd"
	_ "mygf/internal/logic"
	_ "mygf/internal/packed"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
