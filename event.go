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

func newStopEvent() event {
	return event{
		name: "stop",
	}
}

func newRenderEvent() event {
	return event{
		name: "render",
	}
}
