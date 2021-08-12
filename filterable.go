package filterable

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Filterable []interface{}
type orderable Filterable

type emptyFilterableSelection struct{}

var (
	empty = &emptyFilterableSelection{}
)

func New(slice interface{}) (*Filterable, error) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
		return nil, fmt.Errorf("argument not a valid slice")
	}

	size := s.Len()

	filterable := make(Filterable, size)

	for idx := 0; idx < size; idx++ {
		filterable[idx] = s.Index(idx).Interface()
	}

	return &filterable, nil
}

func Empty() *emptyFilterableSelection {
	return empty
}

func Range(start int, count int) *Filterable {
	filterable := Filterable{}

	if count <= 0 {
		return &filterable
	}

	stop := start + count

	for value := start; value < stop; value++ {
		filterable = append(filterable, value)
	}

	return &filterable
}

func (items *Filterable) Unwrap() Filterable {
	return *items
}

func (items *orderable) Unwrap() orderable {
	return *items
}

func (items *orderable) AsFilterable() *Filterable {
	return (*Filterable)(items)
}

func (items *Filterable) AsOrderable() *orderable {
	orderable := orderable{}

	orderable = append(orderable, *items...)
	return &orderable
}

func (items *Filterable) Where(predicate func(interface{}) bool) *Filterable {
	return items.WhereIndexed(func(_ int, key interface{}) bool {
		return predicate(key)
	})
}

func (items *Filterable) WhereIndexed(predicate func(int, interface{}) bool) *Filterable {
	projection := Filterable{}

	for index, item := range *items {
		if predicate(index, item) {
			projection = append(projection, item)
		}
	}

	return &projection
}

func (items *Filterable) Any(predicate func(interface{}) bool) bool {
	for _, item := range *items {
		if predicate(item) {
			return true
		}
	}

	return false
}

func (items *Filterable) All(predicate func(interface{}) bool) bool {
	return !items.Any(func(value interface{}) bool {
		return !predicate(value)
	})
}

func (items *Filterable) Select(keySelector func(interface{}) interface{}) *Filterable {
	return items.SelectIndexed(func(_ int, value interface{}) interface{} {
		return keySelector(value)
	})
}

func (items *Filterable) SelectIndexed(keySelector func(int, interface{}) interface{}) *Filterable {
	projection := Filterable{}

	for index, item := range *items {
		if key := keySelector(index, item); key != empty {
			projection = append(projection, key)
		}
	}

	return &projection
}

func (items *Filterable) Distinct() *Filterable {
	return items.DistinctBy(func(value interface{}) interface{} {
		return value
	})
}

func (items *Filterable) DistinctBy(keySelector func(interface{}) interface{}) *Filterable {
	set := map[interface{}]bool{}

	deduped := Filterable{}

	for _, item := range *items {
		key := keySelector(item)

		if _, seen := set[key]; !seen {
			deduped = append(deduped, item)
			set[key] = true
		}
	}

	return &deduped
}

func (items *Filterable) Union(collection *Filterable) *Filterable {
	projection := append(*items, *collection...)
	return (&projection).Distinct()
}

func (items *Filterable) Intersect(collection *Filterable) *Filterable {
	second := map[interface{}]bool{}

	for _, item := range *collection {
		second[item] = true
	}

	intersection := Filterable{}

	for _, item := range *items {
		if _, exist := second[item]; exist {
			intersection = append(intersection, item)
		}
	}

	return (&intersection).Distinct()
}

func (items *Filterable) Except(collection *Filterable) *Filterable {
	second := map[interface{}]bool{}

	for _, item := range *collection {
		second[item] = true
	}

	projection := Filterable{}

	for _, item := range *items {
		if _, exist := second[item]; !exist {
			projection = append(projection, item)
		}
	}

	return (&projection).Distinct()
}

func (items *Filterable) Skip(count int) *Filterable {
	if items := *items; count >= 0 && len(items) > count {
		projection := items[count:]
		return &projection
	}

	return &Filterable{}
}

func (items *Filterable) SkipWhile(predicate func(interface{}) bool) *Filterable {
	return items.SkipWhileIndexed(func(_ int, value interface{}) bool {
		return predicate(value)
	})
}

func (items *Filterable) SkipWhileIndexed(predicate func(int, interface{}) bool) *Filterable {
	index, size := 0, len(*items)

	for index < size && predicate(index, (*items)[index]) {
		index++
	}

	if index >= size {
		return &Filterable{}
	}

	projection := (*items)[index:]
	return &projection
}

func (items *Filterable) Take(count int) *Filterable {
	if count <= 0 {
		return &Filterable{}
	}

	if items := *items; count > 0 && count < len(items) {
		items = items[:count]
		return &items
	}

	return items
}

func (items *Filterable) TakeWhile(predicate func(interface{}) bool) *Filterable {
	return items.TakeWhileIndexed(func(_ int, value interface{}) bool {
		return predicate(value)
	})
}

func (items *Filterable) TakeWhileIndexed(predicate func(int, interface{}) bool) *Filterable {
	projection := Filterable{}

	for idx := 0; idx < len(*items) && predicate(idx, (*items)[idx]); idx++ {
		projection = append(projection, (*items)[idx])
	}

	return &projection
}

func (items *Filterable) First() interface{} {
	if items := *items; len(items) > 0 {
		return items[0]
	}

	return nil
}

func (items *Filterable) FirstWhere(predicate func(interface{}) bool) interface{} {
	return items.SkipWhile(func(value interface{}) bool {
		return !predicate(value)
	}).First()
}

func (items *Filterable) Last() interface{} {
	if items, size := *items, len(*items); size > 0 {
		return items[size-1]
	}

	return nil
}

func (items *Filterable) LastWhere(predicate func(interface{}) bool) interface{} {
	for items, idx := *items, len(*items)-1; idx >= 0; idx-- {
		if predicate(items[idx]) {
			return items[idx]
		}
	}

	return nil
}

func (items *Filterable) Count() int {
	return len(*items)
}

func (items *Filterable) CountWhere(predicate func(interface{}) bool) interface{} {
	count := 0

	for items, idx, size := *items, 0, len(*items); idx < size; idx++ {
		if predicate(items[idx]) {
			count++
		}
	}

	return count
}

func (items *Filterable) OrderBy(selector func(object interface{}) interface{}) *orderable {
	values := *items.AsOrderable()

	sort.SliceStable(values, func(i, j int) bool {
		first := fmt.Sprintf("%v", selector(values[i]))
		second := fmt.Sprintf("%v", selector(values[j]))

		return first < second
	})

	return &values
}

func (items *Filterable) OrderByDescending(selector func(object interface{}) interface{}) *orderable {
	values := *items.AsOrderable()

	sort.SliceStable(values, func(i, j int) bool {
		first := fmt.Sprintf("%v", selector(values[i]))
		second := fmt.Sprintf("%v", selector(values[j]))

		return first > second
	})

	return &values
}

func (items *Filterable) Order(sortOrder string, selector func(object interface{}) interface{}) *orderable {
	switch strings.ToLower(sortOrder) {
	case "asc":
		return items.OrderBy(selector)
	case "desc":
		return items.OrderByDescending(selector)
	default:
		return (*orderable)(items)
	}
}
