package main

import (
	"fmt"
	"github.com/shaswata56/deepblue/deepblue"
)

func main() {
	var n, scale int
	fmt.Println("User Scale")
	_, _ = fmt.Scanf("%d %d", &n, &scale)
	deepblue.Init(n, scale)
}