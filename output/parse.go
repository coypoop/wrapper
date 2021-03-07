package main

import (
	"fmt"
)

func main() {
	builders := getLogs(3, 14, 9)
	fmt.Printf("builderid %v\n", builders[0].LogId)
}
