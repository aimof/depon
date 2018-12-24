package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aimof/depon"
)

func main() {
	flag.Parse()

	f, err := depon.NewFormatter(".")
	if err != nil {
		log.Fatalln(err)
	}

	switch {
	case *packageFlag != "":
		im, ex, err := f.ShowNode(*packageFlag)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("export: %v\n", ex)
		fmt.Printf("import: %v\n", im)
	default:
		count := f.CountAll()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Ex\tIm\tpackages")
		for name, c := range count {
			fmt.Printf("%d\t%d\t%s\n", c.Parents, c.Children, name)
		}
	}
}
