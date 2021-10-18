package main

import (
	"fmt"
	m "github.com/Polyconseil/k8s-proxy-image-swapper/mutate"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		os.Exit(1)
	}

	for _, v := range os.Args[1:] {
		fmt.Printf("%s ", m.GetPatchedImageUrl(v, "docker.io"))
	}
	fmt.Printf("\n")
}
