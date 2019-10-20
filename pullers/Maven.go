package pullers

type Maven struct {
	Puller
	Name string
}

func (m *Maven) Pull() {
	println("pull " + m.Name)
}
