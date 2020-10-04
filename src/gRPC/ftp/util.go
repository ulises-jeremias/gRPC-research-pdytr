package ftp

// Djb2 function
func Djb2(bytes []byte) uint64 {
	var hash uint64 = 5381

	for _, c := range s {
		hash = ((hash << 5) + hash) + uint64(c)
		// the above line is an optimized version of the following line:
		// hash = hash * 33 + int64(c)
		// which is easier to read and understand...
	}

	return hash
}
