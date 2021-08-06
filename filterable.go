package filterable

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type filterable []interface{}
type orderable filterable

type emptyFilterableSelection struct{}

var (
	empty = &emptyFilterableSelection{}
)

func New(slice interface{}) (*filterable, error) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
		return nil, fmt.Errorf("argument not a valid slice")
	}

	size := s.Len()

	filterable := make(filterable, size)

	for idx := 0; idx < size; idx++ {
		filterable[idx] = s.Index(idx).Interface()
	}

	return &filterable, nil
}

func Empty() *emptyFilterableSelection {
	return empty
}

func Range(start int, count int) *filterable {
	if count <= 0 {
		return &filterable{}
	}

	values, stop, idx := make([]int, count), start+count, 0

	for value := start; value < stop; value++ {
		values[idx] = value
		idx++
	}

	collection, _ := New(values)
	return collection
}

func (items *filterable) Unwrap() filterable {
	return *items
}

func (items *orderable) Unwrap() orderable {
	return *items
}

func (items *orderable) AsFilterable() *filterable {
	return (*filterable)(items)
}

func (items *filterable) AsOrderable() *orderable {
	orderable := orderable{}

	orderable = append(orderable, *items...)
	return &orderable
}

func (items *filterable) Where(predicate func(interface{}) bool) *filterable {
	return items.WhereIndexed(func(_ int, key interface{}) bool {
		return predicate(key)
	})
}

func (items *filterable) WhereIndexed(predicate func(int, interface{}) bool) *filterable {
	projection := filterable{}

	for index, item := range *items {
		if predicate(index, item) {
			projection = append(projection, item)
		}
	}

	return &projection
}

func (items *filterable) Any(predicate func(interface{}) bool) bool {
	for _, item := range *items {
		if predicate(item) {
			return true
		}
	}

	return false
}

func (items *filterable) All(predicate func(interface{}) bool) bool {
	return !items.Any(func(value interface{}) bool {
		return !predicate(value)
	})
}

func (items *filterable) Select(keySelector func(interface{}) interface{}) *filterable {
	return items.SelectIndexed(func(_ int, value interface{}) interface{} {
		return keySelector(value)
	})
}

func (items *filterable) SelectIndexed(keySelector func(int, interface{}) interface{}) *filterable {
	projection := filterable{}

	for index, item := range *items {
		if key := keySelector(index, item); key != empty {
			projection = append(projection, key)
		}
	}

	return &projection
}

func (items *filterable) Distinct() *filterable {
	return items.DistinctBy(func(value interface{}) interface{} {
		return value
	})
}

func (items *filterable) DistinctBy(keySelector func(interface{}) interface{}) *filterable {
	set := map[interface{}]bool{}

	deduped := filterable{}

	for _, item := range *items {
		key := keySelector(item)

		if _, seen := set[key]; !seen {
			deduped = append(deduped, item)
			set[key] = true
		}
	}

	return &deduped
}

func (items *filterable) Union(collection *filterable) *filterable {
	projection := append(*items, *collection...)
	return (&projection).Distinct()
}

func (items *filterable) Intersect(collection *filterable) *filterable {
	second := map[interface{}]bool{}

	for _, item := range *collection {
		second[item] = true
	}

	intersection := filterable{}

	for _, item := range *items {
		if _, exist := second[item]; exist {
			intersection = append(intersection, item)
		}
	}

	return (&intersection).Distinct()
}

func (items *filterable) Except(collection *filterable) *filterable {
	second := map[interface{}]bool{}

	for _, item := range *collection {
		second[item] = true
	}

	projection := filterable{}

	for _, item := range *items {
		if _, exist := second[item]; !exist {
			projection = append(projection, item)
		}
	}

	return (&projection).Distinct()
}

func (items *filterable) Skip(count int) *filterable {
	if items := *items; count >= 0 && len(items) > count {
		projection := items[count:]
		return &projection
	}

	return &filterable{}
}

func (items *filterable) SkipWhile(predicate func(interface{}) bool) *filterable {
	return items.SkipWhileIndexed(func(_ int, value interface{}) bool {
		return predicate(value)
	})
}

func (items *filterable) SkipWhileIndexed(predicate func(int, interface{}) bool) *filterable {
	index, size := 0, len(*items)

	for index < size && predicate(index, (*items)[index]) {
		index++
	}

	if index >= size {
		return &filterable{}
	}

	projection := (*items)[index:]
	return &projection
}

func (items *filterable) Take(count int) *filterable {
	if count <= 0 {
		return &filterable{}
	}

	if items := *items; count > 0 && count < len(items) {
		items = items[:count]
		return &items
	}

	return items
}

func (items *filterable) TakeWhile(predicate func(interface{}) bool) *filterable {
	return items.TakeWhileIndexed(func(_ int, value interface{}) bool {
		return predicate(value)
	})
}

func (items *filterable) TakeWhileIndexed(predicate func(int, interface{}) bool) *filterable {
	projection := filterable{}

	for idx := 0; idx < len(*items) && predicate(idx, (*items)[idx]); idx++ {
		projection = append(projection, (*items)[idx])
	}

	return &projection
}

func (items *filterable) First() interface{} {
	if items := *items; len(items) > 0 {
		return items[0]
	}

	return nil
}

func (items *filterable) FirstWhere(predicate func(interface{}) bool) interface{} {
	return items.SkipWhile(func(value interface{}) bool {
		return !predicate(value)
	}).First()
}

func (items *filterable) Last() interface{} {
	if items, size := *items, len(*items); size > 0 {
		return items[size-1]
	}

	return nil
}

func (items *filterable) LastWhere(predicate func(interface{}) bool) interface{} {
	for items, idx := *items, len(*items)-1; idx >= 0; idx-- {
		if predicate(items[idx]) {
			return items[idx]
		}
	}

	return nil
}

func (items *filterable) Count() int {
	return len(*items)
}

func (items *filterable) CountWhere(predicate func(interface{}) bool) interface{} {
	count := 0

	for items, idx, size := *items, 0, len(*items); idx < size; idx++ {
		if predicate(items[idx]) {
			count++
		}
	}

	return count
}

func (items *filterable) OrderBy(selector func(object interface{}) interface{}) *orderable {
	values := *items.AsOrderable()

	sort.SliceStable(values, func(i, j int) bool {
		first := fmt.Sprintf("%v", selector(values[i]))
		second := fmt.Sprintf("%v", selector(values[j]))

		return first < second
	})

	return &values
}

func (items *filterable) OrderByDescending(selector func(object interface{}) interface{}) *orderable {
	values := *items.AsOrderable()

	sort.SliceStable(values, func(i, j int) bool {
		first := fmt.Sprintf("%v", selector(values[i]))
		second := fmt.Sprintf("%v", selector(values[j]))

		return first > second
	})

	return &values
}

func (items *filterable) Order(sortOrder string, selector func(object interface{}) interface{}) *orderable {
	switch strings.ToLower(sortOrder) {
	case "asc":
		return items.OrderBy(selector)
	case "desc":
		return items.OrderByDescending(selector)
	default:
		return (*orderable)(items)
	}
}
