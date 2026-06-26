package primitives

import "unsafe"

// Checks if a strings data is a subset of another strings data
func IsSubString(str, possibleSubStr string) bool {
	strPtr := uintptr(unsafe.Pointer(unsafe.StringData(str)))
	subStrPtr := uintptr(unsafe.Pointer(unsafe.StringData(possibleSubStr)))

	return strPtr <= subStrPtr && subStrPtr < strPtr + uintptr(len(str))
}
