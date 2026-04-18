package main
import (
"flag"
"fmt"
)
func main() {
d := flag.Bool("debug", false, "dbg")
o := flag.String("overlay", "", "ovl")
flag.Parse()
fmt.Printf("debug: %v, overlay: %q, args: %v\n", *d, *o, flag.Args())
}
