// Package compass is a generic package to determine geo direction and
// work with the directions.
package compass

import "strings"

// Direction represents a geo direction.
type Direction string

const (
	// Nort direction.
	North Direction = "North"
	// South direction.
	South Direction = "South"
	// East direction.
	East Direction = "East"
	// West direction.
	West Direction = "West"
)

// Directions is a list of main compass directions.
var Directions = []Direction{North, East, South, West}

// ReverseDirection returns the reverse direction of d.
func ReverseDirection(d Direction) Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	}
	panic("unreachable")
}

// ParseDirection parses a string value as compass Direction.
// if string value is not a valid geo value ok will be returned
// with a false value.
func ParseDirection(s string) (direction Direction, ok bool) {
	for _, direction := range Directions {
		if strings.ToLower(s) == strings.ToLower(string(direction)) {
			return direction, true
		}
	}
	return "", false
}
