package model

import "fmt"

// Size represents a file size in bytes. It is a distinct uint64 that
// formats as a human-readable string (e.g. "1.23 GB").
type Size uint64

// String formats the size with a binary unit suffix (e.g. "1.23 GB").
func (s Size) String() string {
	if s == 0 {
		return "0 B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	index := 0
	value := float64(s)
	for value >= 1024 && index < len(units)-1 {
		value /= 1024
		index++
	}
	return fmt.Sprintf("%.2f %s", value, units[index])
}

// PercentOf returns s as a percentage of total (e.g. "12.34").
func (s Size) PercentOf(total Size) string {
	if total == 0 {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", float64(s)/float64(total)*100)
}
