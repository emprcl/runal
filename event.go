package runal

type event struct {
	name  string
	value int
}

func newFPSEvent(fps int) event {
	return event{
		name:  "fps",
		value: fps,
	}
}
