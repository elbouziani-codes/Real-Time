package main

import (
	"os"
	"fmt"
	"strings"
)


func main() {
		text := os.Args[1]
		fmt.Println(strings.TrimSpace(text))
}
