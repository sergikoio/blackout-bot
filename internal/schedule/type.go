package schedule

type schedule struct {
	GroupOne   [][]Time `json:"group_1"`
	GroupTwo   [][]Time `json:"group_2"`
	GroupThree [][]Time `json:"group_3"`
	GroupFour  [][]Time `json:"group_4"`
}

type Time struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
