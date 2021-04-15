package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// TODO: restructure to test for -c flag that will execute a single command an
	// d exit; otherwise run the loop

	for {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("🌲 ")
			command, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Printf("❄❅❄❅ Goodbye!!! ❄❅❄❅")
				return
			}

			command = strings.TrimSuffix(command, "\n")
			fmt.Println(command)
		}
	}
}
