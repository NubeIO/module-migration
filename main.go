package main

import (
	"github.com/NubeIO/module-migration/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
