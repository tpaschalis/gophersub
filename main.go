package main

import "fmt"

func main() {
	a, b := ParseWebVTTFile("samples/sample.vtt")
	fmt.Println(a, b)
}
