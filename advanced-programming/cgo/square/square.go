package main

// int square(int x) {
//   return x * x;
// }
import "C"
import "fmt"

func main() {
	res := C.square(4)
	fmt.Println(res)
}
