package main
import (
"fmt"
"os"
)
func main() {
os.Mkdir("testdir", 0755)
defer os.Remove("testdir")
f, err := os.Open("testdir")
if err != nil {
tln("open err:", err)

}
defer f.Close()
info, err := f.Stat()
if err != nil {
tln("stat err:", err)

}
fmt.Printf("IsDir: %v, Mode: %v\n", info.IsDir(), info.Mode())
}
