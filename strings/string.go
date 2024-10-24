package strings

// Limits a string to the predefined size. If the string is bigger than the limit, it will be truncated and the passed suffix will be appended.
//
// Example:
//
//	Limit("this is a long string", 10, "...") // returns "this is a..."
func Limit(str string, limit int, suffix string) string {
	if len(str) <= limit {
		return str
	}
	return str[:limit] + suffix
}
