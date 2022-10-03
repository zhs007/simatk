package battle5

type FuncEachHero func(*Hero)

type FuncAttack func(*Hero, []*Hero, *Battle) (bool, error)
