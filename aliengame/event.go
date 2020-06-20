package aliengame

import "fmt"

// Event is a game event.
type Event interface {
	String() string
}

// sendEvent sends a game event to the listener.
func (w *World) sendEvent(e Event) {
	if w.events != nil {
		w.events <- e
	}
}

// CityDestroyedEvent is emited when a city is removed from the map.
type CityDestroyedEvent struct {
	// City that has been destroyed.
	City *City

	// Aliens on the city that were residing just before the has
	// destroyed.
	Aliens []*Alien
}

func (e CityDestroyedEvent) String() string {
	var alienNames []string
	for _, alien := range e.Aliens {
		alienNames = append(alienNames, alien.Name)
	}
	return fmt.Sprintf("%q has been destroyed by some mad aliens: \n\t%v", e.City.Name, alienNames)
}

// AlienTrappedEvent is emitted when an alien is trapped inside a
// city because the city has no neighbor cities anymore.
type AlienTrappedEvent struct {
	// City that the alien is trapped in.
	City *City

	// Alien that trapped in the city.
	Alien *Alien
}

func (e AlienTrappedEvent) String() string {
	return fmt.Sprintf("alien %q has trapped in city %q", e.Alien.Name, e.City.Name)
}

// CityHasNoNeighboorsEvent is emited when a city has no neighbor
// city around it.
type CityHasNoNeighborsEvent struct {
	City *City
}

func (e CityHasNoNeighborsEvent) String() string {
	return fmt.Sprintf("city %q left with no neighbors", e.City.Name)
}
