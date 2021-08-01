package main

import (
	"fmt"
	"log"

	"github.com/karagulamos/filterable"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7}

	collection, err := filterable.New(numbers)

	if err != nil {
		log.Fatalln(err)
		return
	}

	collection = collection.Where(func(value interface{}) bool {
		return value.(int)%2 == 1
	})

	fmt.Println(collection.Unwrap())
}
