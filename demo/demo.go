// package demo houses the main file for the point location demo.

package main

import (
	"log"

	"github.com/200sc/go-compgeo/demo/demo"

	"github.com/oakmound/oak"
)

func main() {
	// Load our configuration file
	// to initialize the engine
	err := oak.LoadConf("oak.config")
	if err != nil {
		log.Fatal(err)
	}
	// Add some console commands
	demo.AddCommands()
	// Define the only scene in the demo, "demo",
	// which when ended, just resets to itself.
	oak.AddScene("demo", demo.InitScene, demo.LoopScene,
		func() (string, *oak.SceneResult) {
			return "demo", nil
		},
	)
	// Start the engine, beginning at our "demo" scene
	oak.Init("demo")
}
