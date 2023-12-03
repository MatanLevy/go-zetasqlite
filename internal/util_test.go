package internal

import (
	"testing"
)

func TestToLocation(t *testing.T) {
	t.Run("+09", func(t *testing.T) {
		if _, err := toLocation("+09"); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("+09:00", func(t *testing.T) {
		if _, err := toLocation("+09:00"); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("-09", func(t *testing.T) {
		if _, err := toLocation("-09"); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("-09:00", func(t *testing.T) {
		if _, err := toLocation("-09:00"); err != nil {
			t.Fatal(err)
		}
	})
}

func TestSortValues(t *testing.T) {
	t.Run("Order by ascending", func(t *testing.T) {
		orderedValues := []*OrderedValue{
			createOrderedValues([]int64{1, 10, 20}),
			createOrderedValues([]int64{5, 2, 10}),
			createOrderedValues([]int64{5, 2, 1}),
			createOrderedValues([]int64{5, 3, 3}),
			createOrderedValues([]int64{3, 4, 3}),
		}
		orderByList := createOrderByList([]bool{true, true, true})

		sortValues(orderByList, orderedValues)

		expectedValues := [][]int64{
			{1, 10, 20},
			{3, 4, 3},
			{5, 2, 1},
			{5, 2, 10},
			{5, 3, 3},
		}
		validateNewOrder(t, expectedValues, orderedValues)
	})

	t.Run("Order by with ascending and descending", func(t *testing.T) {
		orderedValues := []*OrderedValue{
			createOrderedValues([]int64{1, 10, 20}),
			createOrderedValues([]int64{5, 2, 10}),
			createOrderedValues([]int64{5, 2, 1}),
			createOrderedValues([]int64{5, 3, 3}),
			createOrderedValues([]int64{3, 4, 3}),
		}
		orderByList := createOrderByList([]bool{true, false, false})

		sortValues(orderByList, orderedValues)

		expectedValues := [][]int64{
			{1, 10, 20},
			{3, 4, 3},
			{5, 3, 3},
			{5, 2, 10},
			{5, 2, 1},
		}
		validateNewOrder(t, expectedValues, orderedValues)
	})

	t.Run("No order by values (sort is not performed)", func(t *testing.T) {
		orderedValues := []*OrderedValue{
			createOrderedValues([]int64{1}),
			createOrderedValues([]int64{5}),
			createOrderedValues([]int64{5}),
			createOrderedValues([]int64{5}),
			createOrderedValues([]int64{3}),
		}
		orderByList := createOrderByList([]bool{})

		sortValues(orderByList, orderedValues)

		expectedValues := [][]int64{
			{1},
			{5},
			{5},
			{5},
			{3},
		}
		validateNewOrder(t, expectedValues, orderedValues)
	})
}

// Converts int64 array to a OrderedValue container with the same values (as IntValue)
func createOrderedValues(values []int64) *OrderedValue {
	var orderByObjects []*AggregateOrderBy

	for _, value := range values {
		orderByObjects = append(orderByObjects, &AggregateOrderBy{Value: IntValue(value)})
	}

	return &OrderedValue{OrderBy: orderByObjects}
}

// Converts boolean array, order-by-direction array (whether ascending order or descending order).
// True translates to ascending order. False translates to descending order.
func createOrderByList(values []bool) []*AggregateOrderBy {
	var orderByObjects []*AggregateOrderBy

	for _, value := range values {
		orderByObjects = append(orderByObjects, &AggregateOrderBy{IsAsc: value})
	}

	return orderByObjects
}

// Converts OrderValue container into an array of int64. Assumes all values are of type IntValue
func getIntValues(orderedValues *OrderedValue) []int64 {
	var intValues []int64

	for _, orderedValue := range orderedValues.OrderBy {
		intValue := int64(orderedValue.Value.(IntValue))
		intValues = append(intValues, intValue)
	}

	return intValues
}

// Compares two int64 slices.
// Returns True if lists are equal (all values are equal in the same order). False otherwise.
func slicesEqual(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Validates all values in OrderedValue containers (assumes all values are of type IntValue) match the expected values.
// Fails tests if values don't match.
func validateNewOrder(t *testing.T, expectedValues [][]int64, actualValues []*OrderedValue) {
	for i, expectedValuesArray := range expectedValues {
		actualValuesArray := getIntValues(actualValues[i])
		if !slicesEqual(expectedValuesArray, actualValuesArray) {
			t.Errorf("Sort result is not as expected")
		}
	}
}
