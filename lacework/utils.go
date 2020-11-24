package lacework

// turn an interface array into a string array
func castStringArray(iArray []interface{}) []string {
	return castAndTransformStringArray(iArray, func(s string) string { return s })
}

// turn an interface array into a string array and apply a transformation func
func castAndTransformStringArray(iArray []interface{}, f func(string) string) []string {
	var a []string
	for _, v := range iArray {
		if v == nil {
			continue
		}
		a = append(a, f(v.(string)))
	}
	return a
}

// turn a string array into an instance array
func castStringArrayToInterface(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, str := range strs {
		arr[i] = str
	}
	return arr
}
