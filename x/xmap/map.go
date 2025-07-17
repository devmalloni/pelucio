package xmap

func ToMap[TKey comparable, TValue any](slice []TValue, keyFunc func(TValue) TKey) map[TKey]TValue {
	result := make(map[TKey]TValue, len(slice))
	for _, item := range slice {
		result[keyFunc(item)] = item
	}
	return result
}

func Values[TKey comparable, TValue any](m map[TKey]TValue) []TValue {
	result := make([]TValue, 0, len(m))
	for _, value := range m {
		result = append(result, value)
	}
	return result
}
