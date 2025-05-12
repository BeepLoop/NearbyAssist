package utils

import (
	"time"
)

type Schedule struct {
	Start string
	End   string
}

func (s Schedule) ToTime() (time.Time, time.Time, error) {
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, s.Start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse(layout, s.End)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}

func HasScheduleOverlap(input Schedule, schedules []Schedule) (bool, error) {
	startX, endX, err := input.ToTime()
	if err != nil {
		return false, err
	}

	for _, schedule := range schedules {
		startY, endY, err := schedule.ToTime()
		if err != nil {
			return false, err
		}

		if (startX.Before(endY) || startX.Equal(endY)) && (startY.Before(endX) || startY.Equal(endX)) {
			return true, nil
		}
	}

	return false, nil
}
