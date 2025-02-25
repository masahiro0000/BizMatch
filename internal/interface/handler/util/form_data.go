package util

import "strconv"

// GenerateAgeList returns a slice of integers representing ages from 20 to 80.
func GenerateAgeList() []int {
	ages := make([]int, 0, 40)
	for age := 20; age <= 59; age++ {
		ages = append(ages, age)
	}
	return ages
}

func GenerateGenderList() []string {
	return []string{"男性", "女性"}
}

// ParseNullableInt converts a string to a pointer to int64.
func ParseNullableInt(s string) (*int64, error) {
	if s == "" {
		return nil, nil
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func StringPtr(s string) (*string) {
	return &s
}

func DerefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// ParseNullableInts takes a slice of string values and
// attempt to parse each non-empty string as an int64.
func ParseNullableInts(strs []string) ([]*int64, error) {
	var result []*int64
	for _, s := range strs {
		if s == "" {
			continue
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, &v)
	}
	return result, nil
}

// ParseStringsAsPtrs takes a slice of string values,
// and for each non-empty string.
func ParseStringsAsPtrs(strs []string) []*string {
	var result []*string
	for _, s := range strs {
		if s == "" {
			continue
		}
		copy := s
		result = append(result, &copy)
	}
	return result
}