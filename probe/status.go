
package probe

import (
	"strings"

	"github.com/megaease/easeprobe/global"
)

// Status is the status of Probe
type Status int

// The status of a probe
const (
	StatusInit Status = iota
	StatusUp
	StatusDown
	StatusUnknown
	StatusBad
)

var (
	toTitle = map[Status]string{
		StatusInit:    "Initialization",
		StatusUp:      "Success",
		StatusDown:    "Error",
		StatusUnknown: "Unknown",
		StatusBad:     "Bad",
	}
	toString = map[Status]string{
		StatusInit:    "init",
		StatusUp:      "up",
		StatusDown:    "down",
		StatusUnknown: "unknown",
		StatusBad:     "bad",
	}

	toStatus = global.ReverseMap(toString)

	toEmoji = map[Status]string{
		StatusInit:    "🔎",
		StatusUp:      "✅",
		StatusDown:    "❌",
		StatusUnknown: "⛔️",
		StatusBad:     "🚫",
	}
)

// Title convert the Status to title
func (s Status) Title() string {
	if val, ok := toTitle[s]; ok {
		return val
	}
	return "Unknown"
}

// String convert the Status to string
func (s Status) String() string {
	if val, ok := toString[s]; ok {
		return val
	}
	return "unknown"
}

// Status convert the string to Status
func (s *Status) Status(status string) {
	if val, ok := toStatus[strings.ToLower(status)]; ok {
		*s = val
	} else {
		*s = StatusUnknown
	}
}

// Emoji convert the status to emoji
func (s *Status) Emoji() string {
	if val, ok := toEmoji[*s]; ok {
		return val
	}
	return "⛔️"
}

// UnmarshalYAML is Unmarshal the status
func (s *Status) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return global.EnumUnmarshalYaml(unmarshal, toStatus, s, StatusUnknown, "Status")
}

// MarshalYAML is Marshal the status
func (s Status) MarshalYAML() (interface{}, error) {
	return global.EnumMarshalYaml(toString, s, "Status")
}

// UnmarshalJSON is Unmarshal the status
func (s *Status) UnmarshalJSON(b []byte) (err error) {
	return global.EnumUnmarshalJSON(b, toStatus, s, StatusUnknown, "Status")
}

// MarshalJSON is marshal the status
func (s Status) MarshalJSON() (b []byte, err error) {
	return global.EnumMarshalJSON(toString, s, "Status")
}
