package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/karagulamos/filterable"
)

func main() {
	max := flag.Int("max", 1000, "e.g. primes -max 1000")

	flag.Parse()

	primes := filterable.
		Range(2, *max).
		Where(func(value interface{}) bool {
			number := float64(value.(int))
			return filterable.
				Range(2, int(math.Sqrt(number))-1).
				All(func(divisor interface{}) bool {
					return int(number)%divisor.(int) != 0
				})
		}).Unwrap()

	fmt.Println(primes)
}
