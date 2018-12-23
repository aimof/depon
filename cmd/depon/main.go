package main

import (
	"fmt"
	"log"

	"github.com/aimof/depon"
)

func main() {
	f, err := depon.NewFormatter(".")
	if err != nil {
		log.Fatalln(err)
	}

	count := f.CountAll()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Ex\tIm\tpackages")
	for name, c := range count {
		fmt.Printf("%d\t%d\t%s\n", c.Parents, c.Children, name)
	}
}
