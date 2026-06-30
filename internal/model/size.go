package model

import "fmt"

const bytesPerUnit = 1024

// sizeUnits is the ordered list of binary unit suffixes used when formatting a Size for display.
// PB is the largest unit rendered.
var sizeUnits = [...]string{"B", "KB", "MB", "GB", "TB", "PB"}

// Size is a distinct uint64 that wraps a file size in bytes and formats itself as a human
// readable string, for example "1.23 GB".
type Size uint64

// String formats the receiver as a binary unit value with two decimal places, for example
// "1.23 GB" or "0 B" for zero.
func (s Size) String() string {
	if s == 0 {
		return "0 B"
	}
	index := 0
	value := float64(s)
	for value >= bytesPerUnit && index < len(sizeUnits)-1 {
		value /= bytesPerUnit
		index++
	}
	return fmt.Sprintf("%.2f %s", value, sizeUnits[index])
}

// PercentOf returns the receiver as a percentage of total, in the range 0.0 to 100.0. Returns
// 0 when total is 0 to avoid division by zero and to make "no data" cases render as 0.00%.
func (s Size) PercentOf(total Size) float64 {
	if total == 0 {
		return 0
	}
	return float64(s) / float64(total) * 100
}
