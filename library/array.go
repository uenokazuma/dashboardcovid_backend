package library

func Filter_string(arr []string, cond func(string) bool) []string {

	result := []string{}

	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}

	return result
}
