package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"blackout-bot/internal/servertime"
)

type Schedule struct {
	Sch   [][]Time
	group int
}

func NewSchedule(group int, filename string) (*Schedule, error) {
	sch, err := parse(filename)
	if err != nil {
		return nil, err
	}

	var groupSchedule [][]Time

	switch group {
	case 1:
		groupSchedule = sch.GroupOne
	case 2:
		groupSchedule = sch.GroupTwo
	case 3:
		groupSchedule = sch.GroupThree
	case 4:
		groupSchedule = sch.GroupFour
	default:
		return nil, fmt.Errorf("group is not correct")
	}

	return &Schedule{
		Sch:   groupSchedule,
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
		day = 0
	case time.Tuesday:
		day = 1
	case time.Wednesday:
		day = 2
	case time.Thursday:
		day = 3
	case time.Friday:
		day = 4
	case time.Saturday:
		day = 5
	case time.Sunday:
		day = 6
	}

	return day, hour, nil
}

func NextDay(day int) int {
	if day == 6 {
		return 0
	}
	return day + 1
}

func PreviouslyDay(day int) int {
	if day == 0 {
		return 6
	}
	return day - 1
}

func MinusHour(nowHour int) int {
	if nowHour == 0 {
		return 23
	}
	return nowHour - 1
}

func SchedulesNearby(scheduleFirstEnd, scheduleSecondStart int) bool {
	if scheduleFirstEnd == 24 && scheduleSecondStart == 0 {
		return true
	}

	return false
}

func (s *Schedule) IsScheduleSoon() (bool, Time) {
	now, _ := s.GetScheduleNow()
	if now {
		return false, Time{}
	}

	_, hour, err := GetTimeNow()
	if err != nil {
		return false, Time{}
	}

	timeNow, err := servertime.GetKyivTimeNow()
	if err != nil {
		return false, Time{}
	}
	if timeNow.IsZero() {
		return false, Time{}
	}
	timeNowMinutes := timeNow.Minute()

	nextSchedule, _ := s.GetScheduleNext()
	if MinusHour(nextSchedule.Start) == hour && timeNowMinutes >= 30 {
		return true, nextSchedule
	}

	return false, Time{}
}

func (s *Schedule) GetScheduleNow() (bool, Time) {
	day, hour, err := GetTimeNow()
	if err != nil {
		return false, Time{}
	}

	for _, sch := range s.Sch[day] {
		if sch.Start <= hour && sch.End > hour {
			nextSchedule, nextDay := s.GetScheduleNext()
			prevSchedule, prevDay := s.GetSchedulePreviously()

			if nextDay != day && SchedulesNearby(sch.End, nextSchedule.Start) {
				sch.End = nextSchedule.End
			}
			if prevDay != day && SchedulesNearby(prevSchedule.End, sch.Start) {
				sch.Start = prevSchedule.Start
			}

			return true, sch
		}
	}

	return false, Time{}
}

func (s *Schedule) GetScheduleNext() (time Time, day int) {
	day, hour, err := GetTimeNow()
	if err != nil {
		return Time{}, 0
	}

	for _, sch := range s.Sch[day] {
		if sch.Start > hour {
			return sch, day
		}
	}

	day = NextDay(day)
	for _, sch := range s.Sch[day] {
		return sch, day
	}

	return Time{}, 0
}

func (s *Schedule) GetSchedulePreviously() (time Time, day int) {
	day, hour, err := GetTimeNow()
	if err != nil {
		return Time{}, 0
	}

	index := -1
	for i, sch := range s.Sch[day] {
		if sch.Start < hour {
			index = i
		}
	}
	if index != -1 {
		return s.Sch[day][index], day
	}

	day = PreviouslyDay(day)
	index = len(s.Sch[day])

	return s.Sch[day][index-1], day
}
