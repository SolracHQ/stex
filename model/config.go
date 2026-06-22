package model

import "regexp"

// SortBy describes which field to sort directory contents by.
type SortBy int

const (
	SortByName SortBy = iota
	SortBySize
)

// SortOrder controls ascending or descending sort.
type SortOrder int

const (
	Ascending  SortOrder = iota
	Descending
)

// Grouping controls how files and directories are interleaved in the listing.
type Grouping int

const (
	FilesFirst Grouping = iota
	DirsFirst
	FilesOnly
	DirsOnly
	Mixed
)

// IconStyle controls how file and directory counts are labelled.
type IconStyle int

const (
	IconLetters IconStyle = iota
	IconEmoji
)

// Config holds the current sort, filter, and display settings.
type Config struct {
	SortBy     SortBy
	SortOrder  SortOrder
	Grouping   Grouping
	IconStyle  IconStyle
	ShowHidden bool
	Filter     *regexp.Regexp
}

// DefaultConfig returns size-descending sort with mixed grouping.
func DefaultConfig() Config {
	return Config{
		SortBy:     SortBySize,
		SortOrder:  Descending,
		Grouping:   Mixed,
		IconStyle:  IconLetters,
		ShowHidden: false,
	}
}

// GroupingString returns a short label for a Grouping value.
func GroupingString(g Grouping) string {
	switch g {
	case FilesFirst:
		return "files first"
	case DirsFirst:
		return "dirs first"
	case FilesOnly:
		return "files only"
	case DirsOnly:
		return "dirs only"
	case Mixed:
		return "mixed"
	}
	return ""
}

// Toggle cycles SortBy between SortByName and SortBySize.
func (s *SortBy) Toggle() {
	if *s == SortByName {
		*s = SortBySize
	} else {
		*s = SortByName
	}
}

// Toggle cycles SortOrder between Ascending and Descending.
func (s *SortOrder) Toggle() {
	if *s == Ascending {
		*s = Descending
	} else {
		*s = Ascending
	}
}

// Toggle cycles Grouping through all five modes in order.
func (g *Grouping) Toggle() {
	*g = (*g + 1) % 5
}

// Toggle cycles IconStyle between IconLetters and IconEmoji.
func (is *IconStyle) Toggle() {
	*is = (*is + 1) % 2
}
