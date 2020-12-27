package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	bSize  = 1             // byte size
	kbSize = bSize * 1024  // kilobyte size
	mbSize = kbSize * 1024 // megabyte size
	gbSize = mbSize * 1024 // gigabyte size
	bID    = "B"           // byte id
	kbID   = "Kb"          // kb id
	mbID   = "Mb"          // mb id
	gbID   = "Gb"          // gb id
)

// LabelledBytes2Int64 will convert a string like 10mb, 10gb, 10kb
// to an int64 that represents that number of bytes.
func LabelledBytes2Int64(size string) (int64, error) {

	getValue := func(size string, suffix string, suffixVal int64) (int64, error) {
		strValue := strings.TrimSuffix(size, suffix)

		val, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return 0, err
		}

		return val*suffixVal, nil
	}

	switch {
	case strings.HasSuffix(size, bID):
		return getValue(size, bID, bSize)
	case strings.HasSuffix(size, kbID):
		return getValue(size, kbID, kbSize)
	case strings.HasSuffix(size, mbID):
		return getValue(size, mbID, mbSize)
	case strings.HasSuffix(size, gbID):
		return getValue(size, gbID, gbSize)
	default:
		return 0, fmt.Errorf("'%s' is unrecognised as a binary unit", size)
	}
}
