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
	alienMaxMoveCount = 10000 // TODO unit test this
)

// randIndex used to get a random index number in the [0, length) range.
var randIndex = func(length int) int {
	if length <= 1 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(length)
}

// World is a game world. it consist of cities, roads (directions) and aliens.
type World struct {
	ma sync.Mutex // protects following.
	// mp is game map.
	mp Map
	// aliens are a list of living aliens on the world.
	aliens []*Alien

	// events are used to emit game events when certain game actions happen.
	events chan Event

	// done used to keep track of the status of the world to see if it can
	// be resumed or not.
	done bool
}

// Alien is a living creature came from another world.
type Alien struct {
	// Name is the unique name of the alien.
	Name string

	// CityName is the name of the city that alien is currently residing.
	CityName string

	// MoveCount is the number of times that alien has travelled to another city.
	MoveCount int

	// IsTrapped indicates if alien has trapped inside a city because city does
	// not have any neighbor city left around it.
	IsTrapped bool
}

// New creates a new game world by the given game map. game events sent to the
// events channel but providing it is optional.
func New(mp Map, events chan Event) *World {
	return &World{
		mp:     mp,
		events: events,
	}
}

// SpawnAlien randomly spawns new aliens on the map on different cities.
// it can be used at any time, as much as needed to spawn more aliens on the world.
//
// TODO this can accept ...SpawnAlianOption options to give flexibility on which
// city a given alien should spawn at.
func (w *World) SpawnAlien(count int) {
	w.ma.Lock()
	defer w.ma.Unlock()
	// get an indexable list of city names so they can be randomly picked
	// to place aliens in them.
	var cityNames []string
	for cityName := range w.mp {
		cityNames = append(cityNames, cityName)
	}
	// randomly pick a city for all aliens and send them there.
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

// Resume resumes the game world for one iteration by moving aliens to the neighbor
// cities and making them fight with each other.
// cities and aliens might be destroyed, and directions to the gone cities will be
// removed from existing cities.
//
// certain events will be emited depending on the game actions.
//
// canResume returns with false if all aliens are destroyed or all aliens have
// reached to the max move threshold, in that case the world cannot resume anymore.
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
	// a randomly chosen neighbor city of theirs.
	// thus,
	// - multiple (>=2) aliens may end up in the same city. all aliens in the same
	//   city will fight, all will die and the city will be destroyed.
	// - some cities may have left with no aliens.
	w.moveAliens()
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
			// this mad alien has reached to max move treshold or trapped. the world become
			// to a much better place now!
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

// fightAliens makes the mad aliens in the same city fight which will make them
// all dead. the city and all paths to the city also will be destroyed.
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
			// there are one or no alien on the city, no fight today!
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
