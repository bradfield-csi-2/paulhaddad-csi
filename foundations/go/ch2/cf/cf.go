package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bradfield-csi-2/paulhaddad-csi/foundations/go/ch2/tempconv"
)

// cf converts its numeric arguments to Celsius and Fahrenheit
func main() {
	for _, arg := range os.Args[1:] {
		temp, err := strconv.ParseFloat(arg, 64)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid input temperature: %v\n", temp)
			os.Exit(1)
		}
		fahr := tempconv.Fahrenheit(temp)
		cel := tempconv.Celsius(temp)

		fmt.Printf("%s = %s, %s = %s\n", fahr, tempconv.FToC(fahr), cel, tempconv.CToF(cel))
	}
}
