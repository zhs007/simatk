package battle5

type Scene struct {
	Width  int
	Height int
	Heros  [][]*HeroList
}

func (scene *Scene) AddHero(hero *Hero) {
	if scene.Width == 7 && scene.Height == 3 {
		if hero.TeamIndex == 0 {
			scene.Heros[hero.SY-1][3-hero.SX].AddHero(hero)
		} else {
			scene.Heros[hero.SY-1][3+hero.SX].AddHero(hero)
		}
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
