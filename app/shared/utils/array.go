package utils

// GetStringOrEmpty return value in index or empty string
func GetStringOrEmpty(a []string, index int) string {
	if len(a) >= index+1 {
		return a[index]
	}
	return ""
}

func GetFirstN(n int, text string) string {
	if len(text) > n {
		return text[:n]
	}
	return text
}
