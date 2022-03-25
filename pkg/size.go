package sudogo

type Size struct {
	Width  int
	Height int
}

func (size Size) Area() int {
	return size.Width * size.Height
}
