package battle5

type Scene struct {
	Width  int
	Height int
	Heros  [][]*HeroList
}

func (scene *Scene) AddHero(hero *Hero) {
	if scene.Width == 7 && scene.Height == 3 {
		if hero.TeamIndex == 0 {
			hero.Pos.Y = hero.StaticPos.Y - 1
			hero.Pos.X = 3 - hero.StaticPos.X
		} else {
			hero.Pos.Y = hero.StaticPos.Y - 1
			hero.Pos.X = 3 + hero.StaticPos.X
		}

		scene.Heros[hero.Pos.Y][hero.Pos.X].AddHero(hero)
	}
}

func (scene *Scene) HeroMove(hero *Hero, pos *Pos) {
	scene.Heros[hero.Pos.Y][hero.Pos.X].RemoveHero(hero)

	scene.Heros[pos.Y][pos.X].AddHero(hero)
}

func (scene *Scene) GetHeroNum(pos *Pos) int {
	return scene.Heros[pos.Y][pos.X].GetNum()
}

func (scene *Scene) GetHerosWithPos(pos *Pos) *HeroList {
	return scene.Heros[pos.Y][pos.X]
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
