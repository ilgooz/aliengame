package aliengame

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/ilgooz/aliengame/x/compass"
)

// Map is a game map.
type Map map[string]*City // city name - city pair.

// City is a city in the game map.
type City struct {
	// Name is the unique name of the city.
	Name string

	// HasNoNeighbors shows if city has neighbor cities around it.
	HasNoNeighbors bool

	// Neighbors are the neighbor cities of the city. neighbors has direct walk
	// paths to the city.
	// direction keeps the information about what is the direction of the
	// neighbor relative to the city.
	Neighbors map[compass.Direction]string // direction - city name pair.
}

// re regexp used to parse city & neighboor direction information from each city
// defination line from the map defination -file-.
// city names expected to be made from word chars and optionally dashes.
//
// TODO check `=` token to ensure that only `=` used to bind directions and cities
// together. right now parser does not complain about what token is used.
var cityRe = regexp.MustCompile(`(?m)[\w+-]+`)

// ParseMap parses a map defination by reading from r. it then returns the Map
// representation of the given defination. error is not nil when the defination
// file is syntactically isn't correct or on invalid compass direction.
func ParseMap(r io.Reader) (Map, error) {
	mp := make(Map)
	lr := bufio.NewReader(r)
	var lineNumber int
	for {
		lineNumber++
		line, _, err := lr.ReadLine()
		if err == io.EOF {
			return mp, nil
		}
		if err != nil {
			return mp, err
		}
		words := cityRe.FindAllString(string(line), -1)
		lw := len(words)
		if lw == 0 {
			// allow empty lines in the map defination.
			continue
		}
		if lw < 3 || lw%2 != 1 {
			// found missing city name or, one of the compass direction
			// or city name in one of the direction=city pairs.
			return mp, &CityDefinitionError{lineNumber}
		}
		city := &City{
			Name:      words[0],
			Neighbors: make(map[compass.Direction]string),
		}
		for i := 1; i < lw-1; i += 2 {
			directionStr := words[i]
			cityName := words[i+1]
			direction, ok := compass.ParseDirection(directionStr)
			if !ok {
				return mp, &InvalidDirectionError{lineNumber, directionStr}
			}
			city.Neighbors[direction] = cityName
		}
		mp[city.Name] = city
	}
	return mp, nil
}

// CraftMap makes an analysis on the game map to check map integrity like
// determining impossible neighbor cities (a TODO).
//
// CraftMap decorates the game map:
// - to add cities to the city list that originally do not appear in the city
//   list but exist as a neighbor city of other city.
// - add missing neighbor cities of a city if they are not fully described in
//   the map defination but it is known that they are neighbors after the analysis.
//
// TODO can we guess thourhg other cities and discover more directions?
// TODO check if same compass direction is used multiple times for a city.
func CraftMap(mp Map) error {
	for _, city := range mp {
		// find out the neighbors of the city and check if these neighboors
		// actually present in the map defination. if they don't then add these
		// neighboors to the city list.
		// ensure that each neighbor cities points to each other.
		for direction, neighboorCityName := range city.Neighbors {
			revDirection := compass.ReverseDirection(direction)
			neighboorCity, ok := mp[neighboorCityName]
			if ok {
				// found the neighboor in the city list, make sure the neighboor
				// city back reference to its neighboor.
				neighboorCity.Neighbors[revDirection] = city.Name
				continue
			}
			// could not find the neighboor in the city list, add it.
			newCity := &City{
				Name:      neighboorCityName,
				Neighbors: map[compass.Direction]string{revDirection: city.Name},
			}
			mp[newCity.Name] = newCity
		}
	}
	if len(mp) == 0 {
		return errors.New("there must be at least one city in the map")
	}
	return nil
}

// PrintMap prints a map to w and sorts the cities and directions alphabetically.
func PrintMap(w io.Writer, mp Map) error {
	bw := bufio.NewWriter(w)
	// sort city names.
	var cityNames []string
	for cityName := range mp {
		cityNames = append(cityNames, cityName)
	}
	sort.Strings(cityNames)
	for _, cityName := range cityNames {
		city := mp[cityName]
		bw.WriteString(cityName)
		// sort directions.
		var directions []string
		for direction := range city.Neighbors {
			directions = append(directions, string(direction))
		}
		sort.Strings(directions)
		for _, direction := range directions {
			neighboor := city.Neighbors[compass.Direction(direction)]
			direction := strings.ToLower(direction)
			fmt.Fprintf(bw, " %s=%s", direction, neighboor)
		}
		bw.WriteString("\n")
		if err := bw.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// CityDefinitionError is returned when a city defination in a row is not valid.
type CityDefinitionError struct {
	// LineNumber where error is found.
	LineNumber int
}

func (e *CityDefinitionError) Error() string {
	return fmt.Sprintf("city defination is invalid at line '%d' missing city, neighbor or"+
		"direction in a pair or no direction-neigboor city pair is provided", e.LineNumber)
}

// InvalidDirectionError is returned when direction is not valid.
type InvalidDirectionError struct {
	// LineNumber where error is found.
	LineNumber int

	// Name is the given invalid direction.
	Name string
}

func (e *InvalidDirectionError) Error() string {
	return fmt.Sprintf("invalid direction %q found at line '%d'", e.Name, e.LineNumber)
}
