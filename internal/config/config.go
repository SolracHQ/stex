// Package config holds the user facing configuration model. It is the in memory state every mode
// reads and mutates, the live filter regex, the sort and grouping choices, the display flags.
package config

import "regexp"

// SortBy identifies which field a directory listing is ordered on.
type SortBy int

const (
	SortByName SortBy = iota
	SortBySize
)

// SortOrder identifies the direction of the sort, ascending or descending.
type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

// Grouping identifies how files and subdirectories are interleaved in a listing. The five values
// cover the four block orders and the merged "mixed" order where all items are sorted together by
// the chosen field.
type Grouping int

const (
	FilesFirst Grouping = iota
	DirsFirst
	FilesOnly
	DirsOnly
	Mixed
)

// Config is the full set of user facing settings. In code the fields are read and written
// through methods like Toggle so every state transition goes through one place.
type Config struct {
	SortBy     SortBy         `json:"sort_by"`
	SortOrder  SortOrder      `json:"sort_order"`
	Grouping   Grouping       `json:"grouping"`
	ShowIcons  bool           `json:"show_icons"`
	ShowHidden bool           `json:"show_hidden"`
	LiveFilter bool           `json:"live_filter"`
	Filter     *regexp.Regexp `json:"-"`
}

// DefaultConfig returns the starting state for a first run, largest items first so the user
// immediately sees what is taking the most space.
func DefaultConfig() Config {
	return Config{
		SortBy:     SortBySize,
		SortOrder:  Descending,
		Grouping:   Mixed,
		ShowIcons:  false,
		ShowHidden: false,
		LiveFilter: true,
	}
}

// GroupingString returns the human readable label for a Grouping value, the kind of label that
// fits in the title bar or the settings panel.
func GroupingString(value Grouping) string {
	switch value {
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

// Toggle flips the receiver between SortByName and SortBySize.
func (sort *SortBy) Toggle() {
	if *sort == SortByName {
		*sort = SortBySize
	} else {
		*sort = SortByName
	}
}

// Toggle flips the receiver between Ascending and Descending.
func (order *SortOrder) Toggle() {
	if *order == Ascending {
		*order = Descending
	} else {
		*order = Ascending
	}
}

// Toggle advances the receiver to the next Grouping in declaration order, wrapping back to the
// first value after the last.
func (group *Grouping) Toggle() {
	*group = (*group + 1) % 5
}
