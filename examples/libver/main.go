package main

import (
	"fmt"

	syngo "github.com/synesissoftware/syngo"
	ver2go "github.com/synesissoftware/ver2go"
)

func main() {
	fmt.Printf("syngo v%s\n", syngo.VersionString())
	fmt.Printf("ver2go v%s\n", ver2go.VersionString())
}
