package aliengame

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ilgooz/aliengame/x/compass"
	"github.com/stretchr/testify/require"
)

func TestParseCities(t *testing.T) {
	cases := []struct {
		name        string
		fileContent string
		parsedMap   Map
		parseErr    error
	}{
		{
			"a valid map file",
			`
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
`,
			Map{
				"Foo": &City{
					Name: "Foo",
					Neighbors: map[compass.Direction]string{
						compass.North: "Bar",
						compass.West:  "Baz",
						compass.South: "Qu-ux",
					},
				},
				"Bar": &City{
					Name: "Bar",
					Neighbors: map[compass.Direction]string{
						compass.South: "Foo",
						compass.West:  "Bee",
					},
				},
			},
			nil,
		},
		// TODO test more branches.
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cities, err := ParseMap(strings.NewReader(tt.fileContent))
			require.Equal(t, tt.parseErr, err)
			require.Equal(t, tt.parsedMap, cities)
		})
	}
}

func TestPrintMap(t *testing.T) {
	mapdef := `
Foo north=Bar west=Baz south=Qu-ux
Yee west=Bar
`
	mp, err := ParseMap(strings.NewReader(mapdef))
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, PrintMap(&buf, mp))
	require.Equal(t, "Foo north=Bar south=Qu-ux west=Baz\nYee west=Bar\n", buf.String())
}

func TestCraftMap(t *testing.T) {
	mp := Map{
		"Foo": &City{
			Name: "Foo",
			Neighbors: map[compass.Direction]string{
				compass.North: "Bar",
				compass.South: "Qu-ux",
				compass.West:  "Baz",
			},
		},
		"Bar": &City{
			Name: "Bar",
			Neighbors: map[compass.Direction]string{
				compass.West: "Bee",
			},
		},
	}
	mpcrafted := Map{
		"Bar": &City{
			Name: "Bar",
			Neighbors: map[compass.Direction]string{
				compass.South: "Foo",
				compass.West:  "Bee",
			},
		},
		"Baz": &City{
			Name: "Baz",
			Neighbors: map[compass.Direction]string{
				compass.East: "Foo",
			},
		},
		"Bee": &City{
			Name: "Bee",
			Neighbors: map[compass.Direction]string{
				compass.East: "Bar",
			},
		},
		"Foo": &City{
			Name: "Foo",
			Neighbors: map[compass.Direction]string{
				compass.North: "Bar",
				compass.South: "Qu-ux",
				compass.West:  "Baz",
			},
		},
		"Qu-ux": &City{
			Name: "Qu-ux",
			Neighbors: map[compass.Direction]string{
				compass.North: "Foo",
			},
		},
	}
	mpcraftedjson, _ := json.Marshal(mpcrafted)
	require.NoError(t, CraftMap(mp))
	mpjson, _ := json.Marshal(mp)
	require.Equal(t, mpcraftedjson, mpjson)
}
