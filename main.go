package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("hello world")
	fmt.Printf("is test? %v", strings.HasSuffix(os.Args[0], ".test"))
}
