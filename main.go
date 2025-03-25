package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WIDTH      = 960
	HEIGHT     = 720
	INTER_STAR = 100.0

	MAX_N_PL = 7

	S1_MAS_BEG = 21.0
	S2_MAS_BEG = 7.0
)

func main() {

	var x1 float32 = (S2_MAS_BEG / S1_MAS_BEG * INTER_STAR) / (1 + S2_MAS_BEG/S1_MAS_BEG)
	var x2 float32 = INTER_STAR - x1

	game := Game{
		state:         0,
		star1:         Star{star_mass: S1_MAS_BEG, pos_x: WIDTH / 2, pos_y: HEIGHT/2 + x1, dist: x1},
		star2:         Star{star_mass: S2_MAS_BEG, pos_x: WIDTH / 2, pos_y: HEIGHT/2 - x2, dist: x2},
		addPlanetMass: 5,
		time:          0,
	}

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Double star system")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
