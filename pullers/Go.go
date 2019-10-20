package pullers

type Go struct {
	Puller
}

func (g *Go) Pull() {
	println("pull Go")
}
