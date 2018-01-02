package home

// Interval struct
type Interval struct {
	time int
}

// GetInterval returns the time
func (i Interval) GetInterval() int {
	return i.time
}
