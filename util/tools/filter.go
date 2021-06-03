package tools

func RemoveDuplicateStr(slice []string) []string {
	temp := map[string]struct{}{}
	result := make([]string, 0, len(slice))
	for _, e := range slice {
		l := len(temp)
		temp[e] = struct{}{}
		if len(temp) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func RemoveDuplicateInt(slice []int) []int {
	temp := map[int]struct{}{}
	result := make([]int, 0, len(slice))
	for _, e := range slice {
		l := len(temp)
		temp[e] = struct{}{}
		if len(temp) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
