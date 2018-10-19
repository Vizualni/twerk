package twerk

type Status struct {
	live        int
	working     int
	jobsInQueue int
}

func (s *Status) Idle() int {
	return s.live - s.working
}

func (s *Status) Working() int {
	return s.working
}

func (s *Status) JobsInQueue() int {
	return s.jobsInQueue
}

func (s *Status) Live() int {
	return s.live
}
