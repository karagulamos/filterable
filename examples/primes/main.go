package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/karagulamos/filterable"
)

func main() {
	max := flag.String("max", "1000", "e.g. primes -max 100")

	flag.Parse()

	limit, err := strconv.Atoi(*max)

	if err != nil {
		log.Panicln(err)
		return
	}

	primes := filterable.
		Range(2, limit).
		Where(func(value interface{}) bool {
			number := float64(value.(int))
			return filterable.
				Range(2, int(math.Sqrt(number))-1).
				All(func(divisor interface{}) bool {
					return int(number)%divisor.(int) != 0
				})
		})

	fmt.Println(*primes)
}
