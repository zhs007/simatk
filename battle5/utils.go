package battle5

func Abs[T int | int64](v T) T {
	if v < 0 {
		return -v
	}

	return v
}
