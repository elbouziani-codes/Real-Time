package main 

import (
	"fmt"
)


func test(args ...any) {
	fmt.Println(args)
}
func main() {
	test([]string{"tes", "sh"}...)	
}
