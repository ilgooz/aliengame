package compass

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirectionOrder(t *testing.T) {
	require.Equal(t, North, Directions[0])
	require.Equal(t, East, Directions[1])
	require.Equal(t, South, Directions[2])
	require.Equal(t, West, Directions[3])
}

func TestReverseDirection(t *testing.T) {
	require.Equal(t, North, ReverseDirection(South))
	require.Equal(t, South, ReverseDirection(North))
	require.Equal(t, West, ReverseDirection(East))
	require.Equal(t, East, ReverseDirection(West))
}

func TestParseDirection(t *testing.T) {
	cases := []struct {
		s  string
		d  Direction
		ok bool
	}{
		{"north", North, true},
		{"nOrth", North, true},
		{"North", North, true},
		{"east", East, true},
		{"west", West, true},
		{"south", South, true},
		{"southx", "", false},
	}
	for _, tt := range cases {
		t.Run(tt.s, func(t *testing.T) {
			d, ok := ParseDirection(tt.s)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.d, d)
		})
	}
}
