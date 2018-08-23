package core

// https://github.com/juliangruber/go-intersect/blob/master/intersect.go
// Complexity: O(n^2)
func Intersect(a []string, b []string) []string {
    var result []string
    for _, v := range a {
        if Contains(b, v) {
            result = append(result, v)
        }
    }
	return result
}

func Contains(testArray []string, testValue string) bool {
    var found = false
    for _,v := range testArray {
        if v == testValue {
            found = true
            break
        }
    }
    return found
}
