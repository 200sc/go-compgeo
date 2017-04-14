package main

import (
	"log"

	"github.com/200sc/go-compgeo/demo/demo"

	"bitbucket.org/oakmoundstudio/oak"
)

func main() {
	err := oak.LoadConf("oak.config")
	if err != nil {
		log.Fatal(err)
	}
	demo.AddCommands()
	oak.AddScene("demo", demo.InitScene, demo.LoopScene,
		func() (string, *oak.SceneResult) {
			return "demo", nil
		},
	)
	oak.Init("demo")
}
