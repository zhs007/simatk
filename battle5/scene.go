package battle5

type Scene struct {
	Width  int
	Height int
	Heros  [][]*HeroList
}

func (scene *Scene) AddHero(hero *Hero) {
	if scene.Width == 7 && scene.Height == 3 {
		if hero.TeamIndex == 0 {
			hero.Y = hero.SY - 1
			hero.X = 3 - hero.SX
		} else {
			hero.Y = hero.SY - 1
			hero.X = 3 + hero.SX
		}

		scene.Heros[hero.Y][hero.X].AddHero(hero)
	}
}

func NewScene(w, h int) *Scene {
	scene := &Scene{
		Width:  w,
		Height: h,
	}

	for y := 0; y < h; y++ {
		arr := []*HeroList{}

		for x := 0; x < w; x++ {
			lst := NewHeroList()

			arr = append(arr, lst)
		}

		scene.Heros = append(scene.Heros, arr)
	}

	return scene
}
