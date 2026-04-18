package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	sq, err := NewSquashLayer(os.Args[1])
	if err != nil {
		fmt.Printf("NewSquashLayer err: %v\n", err)
		return
	}
	
	des, err := sq.ReadDir("")
	for _, d := range des {
		name := d.Name()
		childPath := "/" + name
		info, err := sq.Stat(childPath)
		fmt.Printf("Stat(%q) -> err: %v, info: %v\n", childPath, err, info != nil)
	}
}
