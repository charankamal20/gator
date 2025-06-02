package main

import (
	"github.com/charankamal20/gator/internal/config"
)


func main() {
	conf := config.Read()
	conf.SetUser("charankamal20")

	newConf := config.Read()
	newConf.PrintConfig()
}
