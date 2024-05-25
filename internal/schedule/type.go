package schedule

type schedule map[string]map[string]map[string]ElCode

type Time struct {
	Start int
	End   int
	Type  ElCode
}

func (t Time) Reset() Time {
	t.Start = -1
	t.End = -1
	return t
}
