package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"blackout-bot/internal/servertime"
)

type ElCode string

func (e ElCode) ElectricityAvail() bool {
	return e == "yes" || e == "maybe"
}

type Schedule struct {
	Sch   map[int][]Time // map[day]map[hour]
	group int
}

func NewSchedule(group int, filename string) (*Schedule, error) {
	sch, err := parse(filename)
	if err != nil {
		return nil, err
	}

	var groupSchedule map[string]map[string]ElCode

	switch group {
	case 1:
		groupSchedule, _ = sch["1"]
	case 2:
		groupSchedule, _ = sch["2"]
	case 3:
		groupSchedule, _ = sch["3"]
	case 4:
		groupSchedule, _ = sch["4"]
	case 5:
		groupSchedule, _ = sch["5"]
	case 6:
		groupSchedule, _ = sch["6"]
	default:
		return nil, fmt.Errorf("group is not correct")
	}

	return &Schedule{
		Sch:   convertToSchedule(groupSchedule),
		group: group,
	}, nil
}

func parse(filename string) (schedule, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return schedule{}, err
	}

	var s schedule
	err = json.NewDecoder(bytes.NewReader(file)).Decode(&s)
	if err != nil {
		return schedule{}, err
	}

	return s, nil
}

func convertToSchedule(sch map[string]map[string]ElCode) map[int][]Time {
	resp := map[int][]Time{}

	for dayN, day := range sch {
		dayNum, _ := strconv.Atoi(dayN)

		var times []Time
		for hourNum, code := range day {
			hour, _ := strconv.Atoi(hourNum)

			times = append(
				times, Time{
					Start: hour - 1,
					End:   hour,
					Type:  code,
				},
			)
		}

		sort.Slice(
			times, func(i, j int) bool {
				return times[i].Start < times[j].Start
			},
		)
		resp[dayNum] = times
	}

	return resp
}

func GetTimeNow() (day, hour int, err error) {
	timeNow, err := servertime.GetKyivTimeNow()
	if err != nil {
		return 0, 0, err
	}
	if timeNow.IsZero() {
		return 0, 0, fmt.Errorf("time is zero")
	}

	hour = timeNow.Hour()
	switch timeNow.Weekday() {
	case time.Monday:
		day = 1
	case time.Tuesday:
		day = 2
	case time.Wednesday:
		day = 3
	case time.Thursday:
		day = 4
	case time.Friday:
		day = 5
	case time.Saturday:
		day = 6
	case time.Sunday:
		day = 7
	}

	return day, hour, nil
}

func NextDay(day int) int {
	if day == 7 {
		return 1
	}
	return day + 1
}

func (s *Schedule) GetScheduleForDay(day int) []Time {
	var resp []Time

	var currTime = Time{}.Reset()
	for _, t := range s.Sch[day] {
		avail := t.Type.ElectricityAvail()

		if (currTime.Start == -1 && avail) ||
			(currTime.Start != -1 && !avail) {
			continue
		}

		if currTime.Start == -1 && !avail {
			currTime.Start = t.Start
		}
		if currTime.Start != -1 && avail {
			currTime.End = t.Start
			currTime.Type = "no"

			resp = append(resp, currTime)
			currTime = currTime.Reset()
		}
	}

	if currTime.Start != -1 && currTime.End == -1 {
		day := NextDay(day)
		for _, t := range s.Sch[day] {
			avail := t.Type.ElectricityAvail()
			if !avail {
				continue
			} else {
				currTime.End = t.Start
				currTime.Type = "no"
				resp = append(resp, currTime)
				break
			}
		}
	}

	return resp
}
