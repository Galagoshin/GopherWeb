package commands

import "github.com/Galagoshin/GopherWeb/web/framework"

func Stop(string, []string) {
	framework.Shutdown(false)
}
