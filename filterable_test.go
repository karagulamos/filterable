package filterable

import (
	"fmt"
	"reflect"
	"testing"
)

type testScenario struct {
	name     string
	input    interface{}
	expected interface{}
	error    error
	action   func(interface{}) (string, error)
}

var (
	sliceInput   = []int{1, 2, 3, 4, 5, 6, 7}
	arrayInput   = [7]int{1, 2, 3, 4, 5, 6, 7}
	emptyInput   = []int{}
	invalidInput = &sliceInput

	skipCount = 1
	takeCount = 1

	errInvalid = fmt.Errorf("argument not a valid slice")
)

func Test_Filterable_New(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when a valid slice is given",
			input:    sliceInput,
			expected: format_any(sliceInput),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when a valid array is given",
			input:    arrayInput,
			expected: format_any(arrayInput[:]),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when invalid input is given",
			input:    invalidInput,
			expected: format_any(nil),
			error:    errInvalid,
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				return format_any(collection), err
			},
		},
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any([]int{}),
			error:    errInvalid,
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("New", scenarios, t)
}

func Test_Filterable_Range(t *testing.T) {
	type filterableRange struct {
		start int
		stop  int
	}

	scenarios := []testScenario{
		{
			name:     "when a zero count is given",
			input:    &filterableRange{0, 0},
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				r := *input.(*filterableRange)
				collection := Range(r.start, r.stop)
				return format_any(collection.Unwrap()), nil
			},
		},
		{
			name:     "when a negative count is given",
			input:    &filterableRange{1, -7},
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				r := *input.(*filterableRange)
				collection := Range(r.start, r.stop)
				return format_any(collection.Unwrap()), nil
			},
		},
		{
			name:     "when a valid range is given",
			input:    &filterableRange{1, 7},
			expected: format_any(sliceInput),
			action: func(input interface{}) (string, error) {
				r := *input.(*filterableRange)
				collection := Range(r.start, r.stop)
				return format_any(collection.Unwrap()), nil
			},
		},
		{
			name:     "when generating a negative range",
			input:    &filterableRange{-3, 7},
			expected: format_any([]int{-3, -2, -1, 0, 1, 2, 3}),
			action: func(input interface{}) (string, error) {
				r := *input.(*filterableRange)
				collection := Range(r.start, r.stop)
				return format_any(collection.Unwrap()), nil
			},
		},
	}

	run_tests_on("Range", scenarios, t)
}

func Test_Filterable_Any(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(false),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Any(func(value interface{}) bool {
					return value == 4
				})
				return format_any(result), err
			},
		},
		{
			name:     "when value is in slice",
			input:    sliceInput,
			expected: format_any(true),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Any(func(value interface{}) bool {
					return value == 4
				})
				return format_any(result), err
			},
		},
		{
			name:     "when value is not in slice",
			input:    sliceInput,
			expected: format_any(false),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Any(func(value interface{}) bool {
					return value == 8
				})
				return format_any(result), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any(true),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)

				result := collection.
					Select(func(value interface{}) interface{} {
						return value.(int) * 2
					}).
					Any(func(value interface{}) bool {
						return value == 8
					})

				return format_any(result), err
			},
		},
	}

	run_tests_on("Any", scenarios, t)
}

func Test_Filterable_All(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(true),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.All(func(value interface{}) bool {
					return value.(int) > 0 && value.(int) < 8
				})
				return format_any(result), err
			},
		},
		{
			name:     "when all values satisfy predicate",
			input:    sliceInput,
			expected: format_any(true),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.All(func(value interface{}) bool {
					return value.(int) > 0 && value.(int) < 8
				})
				return format_any(result), err
			},
		},
		{
			name:     "when not all values satisfy the predicate",
			input:    sliceInput,
			expected: format_any(false),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.All(func(value interface{}) bool {
					return value.(int) > 0 && value.(int) < 5
				})
				return format_any(result), err
			},
		},
	}

	run_tests_on("All", scenarios, t)
}

func Test_Filterable_Where(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when predicate is satisfied",
			input:    sliceInput,
			expected: format_any([]int{1, 3, 5, 7}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.Where(func(value interface{}) bool {
					return value.(int)%2 == 1
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.Where(func(value interface{}) bool {
					return true
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any([]int{3, 5}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)

				collection = collection.
					Where(func(value interface{}) bool {
						v := value.(int)
						return v > 2 && v < 6
					}).
					Where(func(value interface{}) bool {
						return value.(int)%2 == 1
					})

				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("Where", scenarios, t)
}

func Test_Filterable_WhereIndexed(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when indexed predicate is satisfied",
			input:    sliceInput,
			expected: format_any([]int{1, 3, 5, 7}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.WhereIndexed(func(idx int, value interface{}) bool {
					return idx%2 == 0
				})
				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("WhereIndex", scenarios, t)
}

func Test_Filterable_Select(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.Select(func(value interface{}) interface{} {
					return value
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when valid slice is given",
			input:    sliceInput,
			expected: format_any([]int{2, 4, 6, 8, 10, 12, 14}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Select(func(value interface{}) interface{} {
					return value.(int) * 2
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any([]int{2, 6, 10, 14}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)

				collection = collection.
					Where(func(value interface{}) bool {
						return value.(int)%2 == 1
					}).
					Select(func(value interface{}) interface{} {
						return value.(int) * 2
					})

				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when selector calls Empty() to ignore results",
			input:    sliceInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)

				collection = collection.Select(func(value interface{}) interface{} {
					return Empty()
				})

				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("Select", scenarios, t)
}

func Test_Filterable_SelectIndexed(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when indexed selector is called",
			input:    sliceInput,
			expected: format_any([]int{0, 1, 2, 3, 4, 5, 6}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.SelectIndexed(func(idx int, value interface{}) interface{} {
					return idx
				})
				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("SelectIndexed", scenarios, t)
}

func Test_Filterable_Distinct(t *testing.T) {
	type gameScore struct {
		Name  string
		Score int
	}

	scores := []gameScore{
		{"Alex", 20}, {"Alex", 20},
		{"James", 20}, {"James", 20},
	}

	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				deduped := collection.Distinct()
				return format_any(deduped.Unwrap()), err
			},
		},
		{
			name:     "when slice contains duplicates",
			input:    append(sliceInput, sliceInput...),
			expected: format_any(sliceInput),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				deduped := collection.Distinct()
				return format_any(deduped.Unwrap()), err
			},
		},
		{
			name:     "when deduping a complex type",
			input:    scores,
			expected: format_any([]gameScore{{"Alex", 20}, {"James", 20}}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				deduped := collection.Distinct()
				return format_any(deduped.Unwrap()), err
			},
		},
	}

	run_tests_on("Distinct", scenarios, t)
}

func Test_Filterable_DistinctBy(t *testing.T) {
	type gameScore struct {
		Name  string
		Score int
	}

	scores := []gameScore{
		{"Alex", 20}, {"Alex", 30},
		{"James", 20}, {"James", 30},
	}

	scenarios := []testScenario{
		{
			name:     "when deduping a complex type by key",
			input:    scores,
			expected: format_any([]gameScore{{"Alex", 20}, {"Alex", 30}}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				deduped := collection.DistinctBy(func(value interface{}) interface{} {
					return value.(gameScore).Score
				})
				return format_any(deduped.Unwrap()), err
			},
		},
	}

	run_tests_on("DistinctBy", scenarios, t)
}

func Test_Filterable_Union(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when empty collections are given",
			input:    Range(0, 0),
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				result := Range(0, 0).Union(input.(*filterable))
				return format_any(result.Unwrap()), nil
			},
		},
		{
			name:     "when unioning one non-empty collection with an empty one",
			input:    Range(1, 2),
			expected: format_any([]int{1, 2}),
			action: func(input interface{}) (string, error) {
				result := input.(*filterable).Union(Range(0, 0))
				return format_any(result.Unwrap()), nil
			},
		},
		{
			name:     "when unioning non-empty collections",
			input:    Range(1, 2),
			expected: format_any([]int{1, 2, 3, 4}),
			action: func(input interface{}) (string, error) {
				result := input.(*filterable).Union(Range(3, 2))
				return format_any(result.Unwrap()), nil
			},
		},
		{
			name:     "when unioning dupplicates",
			input:    Range(1, 2),
			expected: format_any([]int{1, 2}),
			action: func(input interface{}) (string, error) {
				result := input.(*filterable).Union(Range(1, 2))
				return format_any(result.Unwrap()), nil
			},
		},
	}

	run_tests_on("Union", scenarios, t)
}

func Test_Filterable_Intersect(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when empty collections are given",
			input:    Range(0, 0),
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				result := Range(0, 0).Intersect(input.(*filterable))
				return format_any(result.Unwrap()), nil
			},
		},
		{
			name:     "when finding intersect of collections with common values",
			input:    Range(0, 7),
			expected: format_any([]int{0, 2, 3}),
			action: func(input interface{}) (string, error) {
				collection, err := New([]int{3, 0, 2})
				result := input.(*filterable).Intersect(collection)
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when finding intersect of collections with duplicates",
			input:    []int{3, 0, 0, 2, 1, 2},
			expected: format_any([]int{3, 0, 2, 1}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Intersect(Range(0, 7))
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when finding intersect of unique collections",
			input:    Range(1, 5),
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				result := input.(*filterable).Intersect(Range(6, 5))
				return format_any(result.Unwrap()), nil
			},
		},
	}

	run_tests_on("Intersect", scenarios, t)
}

func Test_Filterable_Except(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when empty collections are given",
			input:    Range(0, 0),
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				result := Range(0, 0).Except(input.(*filterable))
				return format_any(result.Unwrap()), nil
			},
		},
		{
			name:     "when some values in first collection don't exist in second",
			input:    Range(0, 7),
			expected: format_any([]int{1, 4, 5, 6}),
			action: func(input interface{}) (string, error) {
				collection, err := New([]int{3, 0, 2})
				result := input.(*filterable).Except(collection)
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when all values in first collection don't exist in second",
			input:    Range(1, 5),
			expected: format_any([]int{1, 2, 3, 4, 5}),
			action: func(input interface{}) (string, error) {
				result := input.(*filterable).Except(Range(6, 5))
				return format_any(result.Unwrap()), nil
			},
		},
		{
			name:     "when first collection contains duplicates of values in second",
			input:    []int{3, 0, 0, 2, 1, 2},
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Except(Range(0, 7))
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when first collection contains duplicates of values not in second",
			input:    []int{3, 0, 0, 2, 1, 2},
			expected: format_any([]int{3, 0, 2, 1}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Except(Range(4, 5))
				return format_any(result.Unwrap()), err
			},
		},
	}

	run_tests_on("Except", scenarios, t)
}

func Test_Filterable_Skip(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.Skip(skipCount)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when skip count is negative",
			input:    sliceInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Skip(-skipCount)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when skip count is zero",
			input:    sliceInput,
			expected: format_any(sliceInput),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Skip(0)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when valid slice is given",
			input:    sliceInput,
			expected: format_any([]int{2, 3, 4, 5, 6, 7}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Skip(skipCount)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when skip count greater than length of slice",
			input:    sliceInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Skip(len(sliceInput) + 1)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any([]int{5, 7}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.
					Skip(skipCount).
					Where(func(value interface{}) bool {
						return value.(int)%2 == 1
					}).
					Skip(skipCount)

				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("Skip", scenarios, t)
}

func Test_Filterable_Take(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.Take(takeCount)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when take count is negative",
			input:    sliceInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Take(-takeCount)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when take count is zero",
			input:    sliceInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Take(0)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when valid slice is given",
			input:    sliceInput,
			expected: format_any([]int{1}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.Take(takeCount)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when take count greater than length of slice",
			input:    sliceInput,
			expected: format_any(sliceInput),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.Take(len(sliceInput) + 1)
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any([]int{3}),
			action: func(input interface{}) (string, error) {
				collection, err := New(sliceInput)
				collection = collection.
					Where(func(value interface{}) bool {
						return value.(int)%2 == 1
					}).
					Skip(skipCount).
					Take(takeCount)

				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("Take", scenarios, t)
}

func Test_Filterable_TakeWhile(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.TakeWhile(func(_ interface{}) bool {
					return true
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when a valid slice and predicate is given",
			input:    sliceInput,
			expected: format_any([]int{1}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.TakeWhile(func(value interface{}) bool {
					return value.(int)%2 == 1
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when a truthy predicate is given",
			input:    sliceInput,
			expected: format_any(sliceInput),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.TakeWhile(func(value interface{}) bool {
					return true
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when a falsy predicate is given",
			input:    sliceInput,
			expected: format_any([]int{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.TakeWhile(func(value interface{}) bool {
					return false
				})
				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("TakeWhile", scenarios, t)
}

func Test_Filterable_TakeWhileIndexed(t *testing.T) {
	// https://docs.microsoft.com/en-us/dotnet/api/system.linq.enumerable.takewhile?view=net-5.0
	fruits := []string{
		"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry", "grape", "strawberry",
	}

	expected := []string{
		"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry",
	}

	scenarios := []testScenario{
		{
			name:     "when a valid slice is given",
			input:    fruits,
			expected: format_any(expected),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.TakeWhileIndexed(func(index int, fruit interface{}) bool {
					return len(fruit.(string)) >= index
				})
				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("TakeWhileIndexed", scenarios, t)
}

func Test_Filterable_SkipWhile(t *testing.T) {
	fruits := []string{
		"apple", "passionfruit", "banana", "mango",
		"orange", "blueberry", "grape", "strawberry",
	}

	expected := []string{"grape", "strawberry"}

	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    []string{},
			expected: format_any([]string{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.SkipWhile(func(fruit interface{}) bool {
					return fruit.(string) != "grape"
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when a valid slice is given",
			input:    fruits,
			expected: format_any(expected),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.SkipWhile(func(fruit interface{}) bool {
					return fruit.(string) != "grape"
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when a truthy predicate is given",
			input:    fruits,
			expected: format_any([]string{}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.SkipWhile(func(fruit interface{}) bool {
					return true
				})
				return format_any(collection.Unwrap()), err
			},
		},
		{
			name:     "when a falsy predicate is given",
			input:    fruits,
			expected: format_any(fruits),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.SkipWhile(func(fruit interface{}) bool {
					return false
				})
				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("SkipWhile", scenarios, t)
}

func Test_Filterable_SkipWhileIndexed(t *testing.T) {
	// https://docs.microsoft.com/en-us/dotnet/api/system.linq.enumerable.skipwhile?view=net-5.0
	amounts := []int{
		5000, 2500, 9000, 8000,
		6500, 4000, 1500, 5500,
	}

	expected := []int{4000, 1500, 5500}

	scenarios := []testScenario{
		{
			name:     "when a valid slice is given",
			input:    amounts,
			expected: format_any(expected),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				collection = collection.SkipWhileIndexed(func(index int, amount interface{}) bool {
					return amount.(int) > index*1000
				})
				return format_any(collection.Unwrap()), err
			},
		},
	}

	run_tests_on("SkipWhileIndexed", scenarios, t)
}

func Test_Filterable_First(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(nil),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.First()
				return format_any(item), err
			},
		},
		{
			name:     "when a valid slice is given",
			input:    sliceInput,
			expected: format_any(1),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.First()
				return format_any(item), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any(2),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.Where(func(value interface{}) bool {
					return value.(int)%2 == 0
				}).First()
				return format_any(item), err
			},
		},
	}

	run_tests_on("First", scenarios, t)
}

func Test_Filterable_FirstWhere(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(nil),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.FirstWhere(func(i interface{}) bool {
					return true
				})
				return format_any(item), err
			},
		},
		{
			name:     "when a valid slice is given",
			input:    sliceInput,
			expected: format_any(2),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.FirstWhere(func(value interface{}) bool {
					return value.(int)%2 == 0
				})
				return format_any(item), err
			},
		},
		{
			name:     "when a truthy predicate is given",
			input:    sliceInput,
			expected: format_any(1),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.FirstWhere(func(value interface{}) bool {
					return true
				})
				return format_any(item), err
			},
		},
		{
			name:     "when a falsy predicate is given",
			input:    sliceInput,
			expected: format_any(nil),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.FirstWhere(func(value interface{}) bool {
					return false
				})
				return format_any(item), err
			},
		},
	}

	run_tests_on("FirstWhere", scenarios, t)
}

func Test_Filterable_Last(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(nil),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.Last()
				return format_any(item), err
			},
		},
		{
			name:     "when a valid slice is given",
			input:    sliceInput,
			expected: format_any(7),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.Last()
				return format_any(item), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any(6),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.Where(func(value interface{}) bool {
					return value.(int)%2 == 0
				}).Last()
				return format_any(item), err
			},
		},
	}

	run_tests_on("Last", scenarios, t)
}

func Test_Filterable_LastWhere(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(nil),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.LastWhere(func(i interface{}) bool {
					return true
				})
				return format_any(item), err
			},
		},
		{
			name:     "when a valid slice is given",
			input:    sliceInput,
			expected: format_any(6),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.LastWhere(func(value interface{}) bool {
					return value.(int)%2 == 0
				})
				return format_any(item), err
			},
		},
		{
			name:     "when a truthy predicate is given",
			input:    sliceInput,
			expected: format_any(7),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.LastWhere(func(value interface{}) bool {
					return true
				})
				return format_any(item), err
			},
		},
		{
			name:     "when a falsy predicate is given",
			input:    sliceInput,
			expected: format_any(nil),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.LastWhere(func(value interface{}) bool {
					return false
				})
				return format_any(item), err
			},
		},
	}

	run_tests_on("LastWhere", scenarios, t)
}

func Test_Filterable_Count(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(0),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.Count()
				return format_any(item), err
			},
		},
		{
			name:     "when a valid slice is given",
			input:    sliceInput,
			expected: format_any(7),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.Count()
				return format_any(item), err
			},
		},
		{
			name:     "when chained",
			input:    sliceInput,
			expected: format_any(3),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.Where(func(value interface{}) bool {
					return value.(int)%2 == 0
				}).Count()
				return format_any(item), err
			},
		},
	}

	run_tests_on("Count", scenarios, t)
}

func Test_Filterable_CountWhere(t *testing.T) {
	scenarios := []testScenario{
		{
			name:     "when an empty slice is given",
			input:    emptyInput,
			expected: format_any(0),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.CountWhere(func(value interface{}) bool {
					return value.(int)%2 == 0
				})
				return format_any(result), err
			},
		},
		{
			name:     "when a valid slice is given",
			input:    sliceInput,
			expected: format_any(3),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.CountWhere(func(value interface{}) bool {
					return value.(int)%2 == 0
				})
				return format_any(result), err
			},
		},
		{
			name:     "when a truthy predicate is given",
			input:    sliceInput,
			expected: format_any(7),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.CountWhere(func(value interface{}) bool {
					return true
				})
				return format_any(item), err
			},
		},
		{
			name:     "when a falsy predicate is given",
			input:    sliceInput,
			expected: format_any(0),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				item := collection.CountWhere(func(value interface{}) bool {
					return false
				})
				return format_any(item), err
			},
		},
	}

	run_tests_on("CountWhere", scenarios, t)
}

func Test_Filterable_OrderBy(t *testing.T) {
	sorted := []int{1, 2, 3, 4, 5}
	reversed := []int{5, 4, 3, 2, 1}
	random := []int{2, 1, 3, 5, 4}
	untampered, _ := New(reversed)

	type gameStats struct {
		Name  string
		Score int
	}

	stats := []gameStats{
		{"James", 20}, {"Alex", 30}, {"Alex", 20}, {"James", 30},
	}

	scenarios := []testScenario{
		{
			name:     "when slice is already sorted in ascending",
			input:    sorted,
			expected: format_any(sorted),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderBy(func(object interface{}) interface{} {
					return object.(int)
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when slice is sorted is already sorted in descending",
			input:    reversed,
			expected: format_any(sorted),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderBy(func(object interface{}) interface{} {
					return object
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when input slice is random",
			input:    random,
			expected: format_any(sorted),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderBy(func(object interface{}) interface{} {
					return object
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when sorting in by complex key",
			input:    stats,
			expected: format_any([]gameStats{{"James", 20}, {"Alex", 20}, {"Alex", 30}, {"James", 30}}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderBy(func(object interface{}) interface{} {
					score := object.(gameStats)
					return score.Score
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when filtering an orderable slice",
			input:    stats,
			expected: format_any([]gameStats{{"James", 20}, {"James", 30}}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.
					OrderBy(func(object interface{}) interface{} {
						score := object.(gameStats)
						return score.Score
					}).
					AsFilterable().
					Where(func(item interface{}) bool {
						return item.(gameStats).Name == "James"
					})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when verifying original slice is unchanged",
			input:    untampered,
			expected: format_any(untampered.Unwrap()),
			action: func(input interface{}) (string, error) {
				collection := input.(*filterable)
				collection.OrderBy(func(object interface{}) interface{} {
					return object
				})
				return format_any(collection.Unwrap()), nil
			},
		},
	}

	run_tests_on("OrderBy", scenarios, t)
}

func Test_Filterable_OrderByDescending(t *testing.T) {
	sorted := []int{1, 2, 3, 4, 5}
	reversed := []int{5, 4, 3, 2, 1}
	random := []int{2, 1, 3, 5, 4}
	untampered := Range(1, 10)

	type gameStats struct {
		Name  string
		Score int
	}

	stats := []gameStats{
		{"James", 20}, {"Alex", 30}, {"Alex", 20}, {"James", 30},
	}

	scenarios := []testScenario{
		{
			name:     "when slice is already sorted in descending",
			input:    reversed,
			expected: format_any(reversed),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderByDescending(func(object interface{}) interface{} {
					return object.(int)
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when slice is already sorted in ascending",
			input:    sorted,
			expected: format_any(reversed),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderByDescending(func(object interface{}) interface{} {
					return object
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when input slice is random",
			input:    random,
			expected: format_any(reversed),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderByDescending(func(object interface{}) interface{} {
					return object
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when sorting in by complex key",
			input:    stats,
			expected: format_any([]gameStats{{"Alex", 30}, {"James", 30}, {"James", 20}, {"Alex", 20}}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.OrderByDescending(func(object interface{}) interface{} {
					score := object.(gameStats)
					return score.Score
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when filtering an orderable slice",
			input:    stats,
			expected: format_any([]gameStats{{"James", 30}, {"James", 20}}),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.
					OrderByDescending(func(object interface{}) interface{} {
						score := object.(gameStats)
						return score.Score
					}).
					AsFilterable().
					Where(func(item interface{}) bool {
						return item.(gameStats).Name == "James"
					})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when verifying original slice is unchanged",
			input:    untampered,
			expected: format_any(untampered.Unwrap()),
			action: func(input interface{}) (string, error) {
				collection := input.(*filterable)
				collection.OrderByDescending(func(object interface{}) interface{} {
					return object
				})
				return format_any(collection.Unwrap()), nil
			},
		},
	}

	run_tests_on("OrderByDescending", scenarios, t)
}

func Test_Filterable_Order(t *testing.T) {
	sorted := []int{1, 2, 3, 4, 5}
	reversed := []int{5, 4, 3, 2, 1}

	scenarios := []testScenario{
		{
			name:     "when sorting slice in ascending",
			input:    reversed,
			expected: format_any(sorted),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Order("asc", func(object interface{}) interface{} {
					return object.(int)
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when sorting slice in descending",
			input:    sorted,
			expected: format_any(reversed),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Order("desc", func(object interface{}) interface{} {
					return object
				})
				return format_any(result.Unwrap()), err
			},
		},
		{
			name:     "when invalid sort order is given",
			input:    reversed,
			expected: format_any(reversed),
			action: func(input interface{}) (string, error) {
				collection, err := New(input)
				result := collection.Order("wrong", func(object interface{}) interface{} {
					return object
				})
				return format_any(result.Unwrap()), err
			},
		},
	}

	run_tests_on("Order", scenarios, t)
}

func format_any(collection interface{}) string {
	return fmt.Sprintf("%v", collection)
}

func run_tests_on(method string, scenarios []testScenario, t *testing.T) {
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			got, err := scenario.action(scenario.input)

			if err != nil && err.Error() != scenario.error.Error() {
				t.Errorf("filterable.%v() = %v, expected %v", method, scenario.error, err)
			}

			if !reflect.DeepEqual(got, scenario.expected) {
				t.Errorf("filterable.%v() = %v, expected %v", method, got, scenario.expected)
			}
		})
	}
}
