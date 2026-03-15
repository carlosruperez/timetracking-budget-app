package timer

import (
	"github.com/rupi/timetracking/internal/domain"
)

type Transition struct {
	From   domain.TimerState
	Action string
	To     domain.TimerState
}

var validTransitions = []Transition{
	{domain.TimerStateRunning, "pause", domain.TimerStatePaused},
	{domain.TimerStateRunning, "stop", domain.TimerStateStopped},
	{domain.TimerStatePaused, "resume", domain.TimerStateRunning},
	{domain.TimerStatePaused, "stop", domain.TimerStateStopped},
}

func ValidateTransition(current domain.TimerState, action string) (domain.TimerState, error) {
	for _, t := range validTransitions {
		if t.From == current && t.Action == action {
			return t.To, nil
		}
	}
	return "", domain.ErrInvalidTransition
}
