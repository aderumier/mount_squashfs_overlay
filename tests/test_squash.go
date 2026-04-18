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
fmt.Printf("ReadDir(\"\") err: %v, len: %d\n", err, len(des))
for _, d := range des {
fmt.Printf(" - %s\n", d.Name())
}
}
