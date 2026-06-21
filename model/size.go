package model

import "fmt"

const bytesPerUnit = 1024

var sizeUnits = [...]string{"B", "KB", "MB", "GB", "TB", "PB"}

// Size represents a file size in bytes. It is a distinct uint64 that
// formats as a human-readable string (e.g. "1.23 GB").
type Size uint64

// String formats the size with a binary unit suffix (e.g. "1.23 GB").
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

// PercentOf returns s as a percentage of total (e.g. 12.34).
func (s Size) PercentOf(total Size) float64 {
	if total == 0 {
		return 0
	}
	return float64(s) / float64(total) * 100
}
