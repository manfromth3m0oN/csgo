package util

func RemoveFromSlice(slice []string, s string) []string {
	for i, e := range slice {
		if e == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
