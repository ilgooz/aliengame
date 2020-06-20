// Package aliengame is a game about mad aliens taking over the world.
// They mysteriously appear in various cities, walk to neighboor towns and fight
// with each other because they are so mad.
// Once they start fighting, the whole city teleports to another world in the
// universe and disappears from our big blue bubble.
package aliengame

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ilgooz/aliengame/x/compass"
)

const (
	alienMaxMoveCount = 10000 // unit test this
)

// randIndex used create a random index number as [0, length).
var randIndex = func(length int) int {
	if length <= 1 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(length)
}

// World is a game world. it hosts cities, aliens where aliens fight
// with each other, destroy themselves, cities and paths to cities.
type World struct {
	ma sync.Mutex // protects following.
	// mp is map of the world.
	mp Map
	// aliens are the living aliens on the world.
	aliens []*Alien

	// events are used emit events when certain game action is happened.
	events chan Event

	// done used to keep track of status of the world to see if it can
	// be resumed.
	done bool
}

// Alien is a living creature came from another world.
type Alien struct {
	// Name is the unique name of the alien.
	Name string

	// CityName that alien is current residing.
	CityName string

	// MoveCount is the number of times that alien has travelled to another city.
	MoveCount int

	// IsTrapped indicates if alien has trapped inside a city because city does
	// not have any neighbor city left around it.
	IsTrapped bool
}

// New creates a new game world by given map and send the game events to theevents
// channel. setting the events channel is optional.
func New(mp Map, events chan Event) *World {
	return &World{
		mp:     mp,
		events: events,
	}
}

// SpawnAlien randomly spawns new aliens on the map on different cities.
// it can be at any time and as much as needed to add new aliens.
//
// TODO this can accept ...SpawnAlianOption to provide flexibility on which
// city a given alien should spawn at.
func (w *World) SpawnAlien(count int) {
	w.ma.Lock()
	defer w.ma.Unlock()
	// get an indexable list of city names so they can be randomly picked
	// to place aliens to cities.
	var cityNames []string
	for cityName := range w.mp {
		cityNames = append(cityNames, cityName)
	}
	// randomly pick a city for per alien and place all aliens to a city.
	for i := count; i >= 1; i-- {
		x := randIndex(len(cityNames))
		city := w.mp[cityNames[x]]
		alien := &Alien{
			Name:     fmt.Sprintf("A%d", i),
			CityName: city.Name,
		}
		w.aliens = append(w.aliens, alien)
	}
}

// Resume resumes the game world by one iteration by moving aliens to the neighbor
// cities and making them fight with each other.
// cities and aliens might be destroyed, and paths to gone cities will be removing
// depending the game status.
//
// certain events will be emited depending on the game actions.
//
// canResume returns with false if all aliens are destroyed or all aliens are
// reached to the max move threshold thus, the world cannot resume anymore.
func (w *World) Resume() (canResume bool) {
	w.ma.Lock()
	defer w.ma.Unlock()
	if w.done {
		return false
	}
	defer func() {
		if !canResume {
			w.done = true
			close(w.events)
		}
	}()
	// assuming that every alien that is able to move, will move at the same time to
	// a randomly chosen neighboor city of theirs.
	// thus,
	// - multiple (>=2) aliens may end up in the same city. all aliens in the same
	//   city will fight, all will die and the city will be destroyed.
	// - some cities may have left with no aliens.
	w.moveAliens()
	for _, a := range w.aliens {

		fmt.Println(a.CityName)
	}
	w.fightAliens()
	// check if the world can be resumed again.
	for _, alien := range w.aliens {
		if canAlienMove(alien) {
			return true
		}
	}
	return false
}

func canAlienMove(alien *Alien) bool {
	return alien.MoveCount < alienMaxMoveCount && !alien.IsTrapped
}

// moveAliens moves aliens to the neighbor cities if possible.
func (w *World) moveAliens() {
	for _, alien := range w.aliens {
		if !canAlienMove(alien) {
			// this mad alien has reached to max move treshold. the world become
			// to a much better place now.
			continue
		}
		alien.MoveCount++
		city := w.mp[alien.CityName]
		// list all directions/neighbors.
		var directions []compass.Direction
		for direction := range city.Neighbors {
			directions = append(directions, direction)
		}
		ld := len(directions)
		if ld == 0 {
			// a mad alien has trapped to a city.
			alien.IsTrapped = true
			w.sendEvent(AlienTrappedEvent{
				City:  city,
				Alien: alien,
			})
			continue
		}
		// randomly pick a neighbor and send alien to that city.
		chosenDirection := directions[randIndex(ld)]
		alien.CityName = city.Neighbors[chosenDirection]
	}
}

// fightAliens makes mad aliens in the same city fight resulting them being dead
// and city and all paths to that it being destroyed.
func (w *World) fightAliens() {
	for _, city := range w.mp {
		// find the aliens residing on the city.
		var aliens []*Alien
		for _, alien := range w.aliens {
			if alien.CityName == city.Name {
				aliens = append(aliens, alien)
			}
		}
		if len(aliens) < 2 {
			// there are one or no alien on the city, no alien to fight.
			continue
		}
		// ops! >2 aliens are in the city, they fought!
		// now delete the aliens and city.
		delete(w.mp, city.Name)
		for i := len(w.aliens) - 1; i >= 0; i-- {
			if w.aliens[i].CityName == city.Name {
				w.aliens = append(w.aliens[:i], w.aliens[i+1:]...)
			}
		}
		w.sendEvent(CityDestroyedEvent{
			City:   city,
			Aliens: aliens,
		})
	}
	for _, city := range w.mp {
		// remove danling neigboors (the cities that are no longer exist in the map but
		// referenced by the existing cities).
		for direction, neighboorCityName := range city.Neighbors {
			if _, ok := w.mp[neighboorCityName]; !ok {
				delete(city.Neighbors, direction)
			}
		}
		// send no neighboors left event, once, if a city left out with no
		// neighboors.
		if len(city.Neighbors) == 0 && !city.HasNoNeighbors {
			city.HasNoNeighbors = true
			w.sendEvent(CityHasNoNeighborsEvent{
				City: city,
			})
		}
	}
}

// Map gets a snapshot of current status of the city Map.
func (w *World) Map() Map {
	w.ma.Lock()
	defer w.ma.Unlock()
	return w.mp
}
