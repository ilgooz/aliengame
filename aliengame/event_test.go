package aliengame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// make sure events implements Event.
var _ Event = (*CityDestroyedEvent)(nil)
var _ Event = (*CityHasNoNeighborsEvent)(nil)
var _ Event = (*AlienTrappedEvent)(nil)

func TestCityDestroyedEvent(t *testing.T) {
	require.Equal(t, "\"1\" has been destroyed by some mad aliens: \n\t[2 3]", CityDestroyedEvent{
		City: &City{
			Name: "1",
		},
		Aliens: []*Alien{
			{
				Name: "2",
			},
			{
				Name: "3",
			},
		},
	}.String())
}

func TestCityHasNoNeighbors(t *testing.T) {
	require.Equal(t, "city \"1\" left with no neighbors", CityHasNoNeighborsEvent{
		City: &City{
			Name: "1",
		},
	}.String())
}

func TestAlienTrappedEvent(t *testing.T) {
	require.Equal(t, "alien \"2\" has trapped in city \"1\"", AlienTrappedEvent{
		City: &City{
			Name: "1",
		},
		Alien: &Alien{
			Name: "2",
		},
	}.String())
}
