package main

import (
	"fmt"
	"os"
)

func extract(arg string) (ret []string) {
	start := 0
	depth := 0
	for i, c := range arg {
		if c == '<' {
			if depth == 0 {
				start = i + 1
			}
			depth++
		} else if c == '>' {
			if depth > 0 {
				depth--
				if depth == 0 {
					ret = append(ret, arg[start:i])
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unbalanced braces\n")
			}
		}
	}
	return ret
}

func main() {
	for _, arg := range os.Args[1:] {
		ex := extract(arg)
		fmt.Printf("%q: ", arg)
		for _, p := range ex {
			fmt.Printf("%q ", p)
		}
		fmt.Println()
	}
}
