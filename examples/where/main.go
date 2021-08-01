package main

import (
	"log"

	"github.com/karagulamos/filterable"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7}

	collection, err := filterable.New(numbers)

	if err != nil {
		log.Fatalln(err)
	}

	collection = collection.Where(func(value interface{}) bool {
		return value.(int)%2 == 1
	})

	log.Println(*collection...)
}
