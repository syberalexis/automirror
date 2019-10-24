package pullers

type Puller interface {
	Pull() int
}
