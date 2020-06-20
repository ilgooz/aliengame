package aliengame

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO add tests to test fixed behavior of game engine by elimating randomness.
// in order to do that replace randIndex() and its return value to test different
// behaviors as needed.

func TestGame(t *testing.T) {
	mapdef := `
Foo north=Bar west=Baz south=Qu-ux
Bee south=Bar
Yee west=Bar
`
	// run the game test multiple times to ensure randomnees does not
	// hide bugs in mysterious ways.
	for i := 0; i < 100; i++ {
		i := i
		t.Run(fmt.Sprintf("T%d", i), func(t *testing.T) {
			t.Parallel()
			testGame(t, mapdef)
		})
	}
}

func testGame(t *testing.T, mapdef string) {
	var events []Event
	eventC := make(chan Event)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range eventC {
			events = append(events, event)
		}
	}()
	mp, err := ParseMap(strings.NewReader(mapdef))
	require.NoError(t, err)
	require.NoError(t, CraftMap(mp))
	world := New(mp, eventC)
	world.SpawnAlien(3)
	for world.Resume() {
	}
	wg.Wait()

	mps := world.Map()
	for _, event := range events {
		switch e := event.(type) {
		case CityDestroyedEvent:
			for _, city := range mps {
				require.NotEqual(t, e.City, city)
			}
		case CityHasNoNeighborsEvent:
			if city, ok := mps[e.City.Name]; ok {
				require.Empty(t, city.Neighbors)
			}
		}
	}
}

func TestRandIndex(t *testing.T) {
	cases := []struct {
		length int
		lte    int
	}{
		{10, 9},
		{3, 2},
		{0, 0},
		{1, 0},
		{-1, 0},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			val := randIndex(tt.length)
			require.True(t, val >= 0)
			require.True(t, val <= tt.lte)
		})
	}
}
