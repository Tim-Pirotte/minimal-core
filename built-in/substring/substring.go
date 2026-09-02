package substring

import "unsafe"

func GetStringPtr(str string) uintptr {
	return uintptr(unsafe.Pointer(unsafe.StringData(str)))
}

// Checks if a strings data is a subset of another strings data
func IsSubString(str, possibleSubStr string) bool {
	strPtr := GetStringPtr(str)
	subStrPtr := GetStringPtr(possibleSubStr)

	return strPtr <= subStrPtr && subStrPtr < strPtr + uintptr(len(str))
}
