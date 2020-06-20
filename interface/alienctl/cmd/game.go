package aliengamecmd

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ilgooz/aliengame/aliengame"
	"github.com/spf13/cobra"
)

var (
	mapFilePath string
	alienCount  int
)

// New returns a new alienctl command that can be attached to a cli app.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alienctl",
		Short: "fight aliens, destroy cities!",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handler(mapFilePath, alienCount, cmd.OutOrStdout())
		},
	}
	cmd.Flags().StringVarP(&mapFilePath, "map-file", "m", "", "path to the map file (required)")
	cmd.Flags().IntVarP(&alienCount, "alien-count", "a", 0, "number of aliens to spawn (required)")
	cmd.MarkFlagRequired("map-file")
	cmd.MarkFlagRequired("alien-count")
	return cmd
}

// handler runs the game by using given inputs.
func handler(mapFilePath string, alienCount int, w io.Writer) error {
	mapFile, err := os.Open(mapFilePath)
	if err != nil {
		return err
	}
	defer mapFile.Close()

	// start the game and print game events.
	events := make(chan aliengame.Event)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range events {
			fmt.Fprintf(w, "e>%s\n", event)
		}
	}()
	mp, err := aliengame.ParseMap(mapFile)
	if err != nil {
		return err
	}
	if err := aliengame.CraftMap(mp); err != nil {
		return err
	}
	world := aliengame.New(mp, events)
	world.SpawnAlien(alienCount)
	for world.Resume() {
	}
	wg.Wait()

	// print map state.
	fmt.Fprint(w, "\nMAP STATE:\n")
	return aliengame.PrintMap(w, world.Map())
}
