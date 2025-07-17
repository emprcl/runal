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

func newStartEvent() event {
	return event{
		name: "start",
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

func newExitEvent() event {
	return event{
		name: "exit",
	}
}
