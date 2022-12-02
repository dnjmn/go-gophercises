package main

import (
	"fmt"
	"gophercises/link"
	"log"
	"os"
	"strings"
)

func main() {
	b, e := os.ReadFile("ex.html")
	if e != nil {
		log.Fatal("main: ", e)
	}

	r := strings.NewReader(string(b))
	l, e := link.Parse(r)
	if e != nil {
		log.Fatal("main: ", e)
	}
	fmt.Println(l)
}
