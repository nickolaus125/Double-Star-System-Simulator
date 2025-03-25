package main

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	state int // 0 - menu, 1 - in program, 2-plot

	planet        []Planet
	star1         Star
	star2         Star
	barycenter_x  float32
	barycenter_y  float32
	addPlanetMass int
	time          float32

	star1_x  []float32
	star1_y  []float32
	star2_x  []float32
	star2_y  []float32
	planet_x [MAX_N_PL][]float32
	planet_y [MAX_N_PL][]float32
}

func (g *Game) Update() error {

	g.getInput()

	if g.addPlanetMass < 0 {
		g.addPlanetMass = 0
	}

	g.barycenter_x = (g.star1.pos_x*g.star1.star_mass + g.star2.pos_x*g.star2.star_mass) / (g.star1.star_mass + g.star2.star_mass)

	g.barycenter_y = (g.star1.pos_y*g.star1.star_mass + g.star2.pos_y*g.star2.star_mass) / (g.star1.star_mass + g.star2.star_mass)

	if g.state == 1 {
		g.time += 0.01
		if g.time >= 36 {
			g.time = 0
		}

		g.star1.pos_x, g.star1.pos_y = g.star1.CalcStarPos(g.time, g.barycenter_x, g.barycenter_y, g.star1.star_mass, g.star2.star_mass)
		g.star2.pos_x, g.star2.pos_y = g.star2.CalcStarPos(g.time, g.barycenter_x, g.barycenter_y, g.star1.star_mass, g.star2.star_mass)

		for i := 0; i < len(g.planet); i++ {
			g.planet[i].pos_x, g.planet[i].pos_y = g.planet[i].CalcPlanetPos(g.time, g.star1)
			g.planet[i].pos_x, g.planet[i].pos_y = g.planet[i].CalcPlanetPos(g.time, g.star2)
			for j := 0; j < len(g.planet); j++ {
				if j != i {
					g.planet[i].pos_x, g.planet[i].pos_y = g.planet[i].PlanetToPlanet(g.time, g.planet[j])
				}
			}
		}
		g.star1_x = append(g.star1_x, g.star1.pos_x)
		g.star1_y = append(g.star1_y, g.star1.pos_y)
		g.star2_x = append(g.star2_x, g.star2.pos_x)
		g.star2_y = append(g.star2_y, g.star2.pos_y)

		for k := 0; k < len(g.planet); k++ {
			g.planet_x[k] = append(g.planet_x[k], g.planet[k].pos_x)
			g.planet_y[k] = append(g.planet_y[k], g.planet[k].pos_y)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 20, 68, 255})

	g.writeStrs(screen, g.star1.star_mass, g.star2.star_mass)

	g.DrawStar(screen, g.star1)
	g.DrawStar(screen, g.star2)

	for i := 0; i < len(g.planet); i++ {
		g.DrawPlanet(screen, g.planet[i])
	}

	if g.state == 2 {
		vector.DrawFilledRect(screen, 0, 0, WIDTH, HEIGHT, color.RGBA{255, 255, 255, 255}, false)
		for k := 0; k < len(g.star1_x); k++ {
			vector.DrawFilledCircle(screen, g.star1_x[k], g.star1_y[k], 2, color.RGBA{255, 0, 0, 255}, false)
			vector.DrawFilledCircle(screen, g.star2_x[k], g.star2_y[k], 2, color.RGBA{231, 255, 0, 255}, false)
		}
		for i := 0; i < len(g.planet); i++ {
			for j := 0; j < len(g.planet_x[i]); j++ {
				vector.DrawFilledCircle(screen, g.planet_x[i][j], g.planet_y[i][j], 2, color.RGBA{0, 0, 0, 255}, false)
			}
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *Game) getInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if g.state == 1 {
			g.state = 0
		} else {
			g.state = 1
		}
	}
	if g.state == 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.addPlanetMass++
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.addPlanetMass--
		} else if inpututil.IsKeyJustPressed(ebiten.Key1) {
			g.star1.star_mass++
		} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
			g.star1.star_mass--
		} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
			g.star2.star_mass++
		} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
			g.star2.star_mass--
		} else if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			g.RemovePlanet()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			clear(g.star1_x)
			clear(g.star1_y)
			clear(g.star2_x)
			clear(g.star2_y)
			for i := 0; i < len(g.planet); i++ {
				clear(g.planet_x[i])
				clear(g.planet_y[i])
			}

		} else if inpututil.IsKeyJustPressed(ebiten.KeyT) {
			g.state = 2
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			a, b := ebiten.CursorPosition()
			if len(g.planet) < MAX_N_PL {
				if g.canAddPlanet(float32(a), float32(b)) {
					g.planet = append(g.planet, Planet{planet_mass: float32(g.addPlanetMass), pos_x: float32(a), pos_y: float32(b), dist: 100})
				}
			}
		}
	} else if g.state == 2 {
		if inpututil.IsKeyJustPressed(ebiten.KeyT) {
			g.state = 0
		}
	}
}

func (g *Game) writeStrs(screen *ebiten.Image, mass1 float32, mass2 float32) {
	msg := "Star1: \n" + "    Mass: " + strconv.FormatFloat(float64(mass1), 'f', -1, 32) + "\n Star2: \n" + "    Mass: " + strconv.FormatFloat(float64(g.star2.star_mass), 'f', -1, 32) + "\n  Barycenter x: " + strconv.FormatFloat(float64(g.barycenter_x), 'f', -1, 32) + "\n  Barycenter y: " + strconv.FormatFloat(float64(g.barycenter_y), 'f', -1, 32) + "\n Add Planet of Mass: " + strconv.FormatFloat(float64(g.addPlanetMass), 'f', -1, 32)
	text.Draw(screen, msg, basicfont.Face7x13, 50, 50, color.White)
}

func (g *Game) DrawPlanet(screen *ebiten.Image, planet1 Planet) {
	vector.DrawFilledCircle(screen, planet1.pos_x, planet1.pos_y, 5, color.RGBA{150, 200, 200, 255}, false)
}

func (g *Game) DrawStar(screen *ebiten.Image, star1 Star) {
	vector.DrawFilledCircle(screen, star1.pos_x, star1.pos_y, min(star1.star_mass, 20), color.RGBA{200 - 10, 100, 100, 255}, false)
}

func (g *Game) canAddPlanet(x float32, y float32) bool {
	dist := math.Sqrt(float64((g.star1.pos_x-x)*(g.star1.pos_x-x) + (g.star1.pos_y-y)*(g.star1.pos_y-y)))

	dist2 := math.Sqrt(float64((g.star2.pos_x-x)*(g.star2.pos_x-x) + (g.star2.pos_y-y)*(g.star2.pos_y-y)))

	if dist < 25 || dist2 < 25 {
		return false
	}

	for i := 0; i < len(g.planet); i++ {
		dist = math.Sqrt(float64((g.planet[i].pos_x-x)*(g.planet[i].pos_x-x) + (g.planet[i].pos_y-y)*(g.planet[i].pos_y-y)))
		if dist < 8 {
			return false
		}
	}
	return true
}

func (g *Game) RemovePlanet() {
	if len(g.planet) > 0 {
		g.planet = g.planet[:len(g.planet)-1]
	}
}
