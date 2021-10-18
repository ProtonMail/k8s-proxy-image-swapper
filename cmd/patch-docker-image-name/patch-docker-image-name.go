package main

import (
	"fmt"
	m "github.com/Polyconseil/k8s-proxy-image-swapper/mutate"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	registry := os.Args[1]

	for _, v := range os.Args[2:] {
		fmt.Printf("%s ", m.GetPatchedImageUrl(v, registry))
	}
	fmt.Printf("\n")
}
