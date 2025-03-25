package main

import "math"

type Planet struct {
	planet_id   int
	planet_mass float32
	pos_x       float32
	pos_y       float32
	dist        float32
}

type Star struct {
	star_mass float32
	pos_x     float32
	pos_y     float32
	dist      float32
}

func (star11 Star) CalcStarPos(t float32, b_x float32, b_y float32, mass1 float32, mass2 float32) (float32, float32) {

	podpierw := float64((star11.pos_x-b_x)*(star11.pos_x-b_x) + (star11.pos_y-b_y)*(star11.pos_y-b_y))

	star11.dist = float32(math.Sqrt(podpierw))
	sinb := float32((float64(star11.pos_x - b_x))) / -star11.dist

	cosb := float32((float64(star11.pos_y - b_y))) / star11.dist

	var v_x float32
	var v_y float32

	var v float32 = float32(math.Sqrt(float64(mass1+mass2-star11.star_mass) / float64(INTER_STAR)))

	v_y = sinb * v
	v_x = cosb * v

	return star11.pos_x + v_x, star11.pos_y + v_y

}

func (planet11 Planet) CalcPlanetPos(t float32, star11 Star) (float32, float32) {

	b_x := star11.pos_x
	b_y := star11.pos_y

	podpierw := float64((planet11.pos_x-b_x)*(planet11.pos_x-b_x) + (planet11.pos_y-b_y)*(planet11.pos_y-b_y))

	planet11.dist = float32(math.Sqrt(podpierw))
	sinb := float32((float64(planet11.pos_x - b_x))) / -planet11.dist

	cosb := float32((float64(planet11.pos_y - b_y))) / planet11.dist

	if (star11.pos_x-b_x) > 0 && (star11.pos_y-b_y) > 0 {
		sinb *= -1
	}

	var v float32 = float32(math.Sqrt(float64(star11.star_mass) / float64(planet11.dist)))

	var v_x float32 = 10 * cosb / (planet11.dist) * planet11.planet_mass
	var v_y float32 = 10 * sinb / (planet11.dist) * planet11.planet_mass

	v_x = cosb * v
	v_y = sinb * v

	return planet11.pos_x + v_x, planet11.pos_y + v_y
}

func (planet11 Planet) PlanetToPlanet(t float32, planet22 Planet) (float32, float32) {
	b_x := (planet11.pos_x + planet22.pos_x) / 2
	b_y := (planet11.pos_y + planet22.pos_y) / 2

	podpierw := float64((planet11.pos_x-b_x)*(planet11.pos_x-b_x) + (planet11.pos_y-b_y)*(planet11.pos_y-b_y))

	planet11.dist = float32(math.Sqrt(podpierw))
	sinb := float32((float64(planet11.pos_x - b_x))) / -planet11.dist

	cosb := float32((float64(planet11.pos_y - b_y))) / planet11.dist

	if (planet11.pos_x-b_x) > 0 && (planet11.pos_y-b_y) > 0 {
		sinb *= -1
	}

	var v float32 = float32(math.Sqrt(float64(0.1*planet22.planet_mass) / float64(planet11.dist)))

	var v_x float32 = 1 * cosb / (planet11.dist) * planet11.planet_mass
	var v_y float32 = 1 * sinb / (planet11.dist) * planet11.planet_mass

	v_x = cosb * v
	v_y = sinb * v

	return planet11.pos_x + v_x, planet11.pos_y + v_y
}
