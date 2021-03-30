package tempconv

import "fmt"

type Fahrenheit float64
type Celsius float64

const FreezingC = 32

func CToF(temp Celsius) Fahrenheit {
	return Fahrenheit(temp*9/5 + FreezingC)
}

func FToC(temp Fahrenheit) Celsius {
	return Celsius((temp - FreezingC) * 5 / 9)
}

func (f Fahrenheit) String() string {
	return fmt.Sprintf("%.2f°F", f)
}

func (c Celsius) String() string {
	return fmt.Sprintf("%.2f°C", c)
}
