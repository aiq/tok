package tok

func midSubString(str string, from int) string {
	c := 0
	for i := range str {
		if c == from {
			return str[i:]
		}
		c++
	}
	return ""
}

func subStringFrom(str string, from int, n int) string {
	tail := midSubString(str, from)
	c := 0
	for i := range tail {
		if c == n {
			return str[from : from+i]
		}
		c++
	}
	return tail
}
