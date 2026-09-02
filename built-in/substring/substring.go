package substring

import "unsafe"

func GetStringPtr(str string) uintptr {
    return uintptr(unsafe.Pointer(unsafe.StringData(str)))
}

// Checks if a strings data is a subset of another strings data
func IsSubString(str, possibleSubStr string) bool {
    strPtr := GetStringPtr(str)
    subStrPtr := GetStringPtr(possibleSubStr)

    strEnd := strPtr + uintptr(len(str))
	subStrEnd := subStrPtr + uintptr(len(possibleSubStr))

    return strPtr <= subStrPtr && subStrEnd <= strEnd
}
