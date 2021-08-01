package main

import (
	"fmt"

	"github.com/karagulamos/filterable"
)

func main() {
	pageNumber, pageSize := 3, 20

	values := filterable.
		Range(1, 100).
		Skip((pageNumber - 1) * pageSize).
		Take(pageSize)

	fmt.Println(*values)
}
