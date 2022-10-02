package battle5

type Scene struct {
	Width  int
	Height int
	Heros  [][]*HeroList
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
