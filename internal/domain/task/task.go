package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	Daily        RecurrenceType = "daily"
	Monthly      RecurrenceType = "monthly"
	SpecificDate RecurrenceType = "specific_date"
	EvenDays     RecurrenceType = "even_days"
	OddDays      RecurrenceType = "odd_days"
)

type Task struct {
	ID              int64          `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Status          Status         `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	RecurrenceType  RecurrenceType `json:"recurrence_type,omitempty"`
	RecurrenceValue int            `json:"recurrence_value,omitempty"`
	SpecificDates   []time.Time    `json:"specific_dates,omitempty"`
	EndDate         time.Time      `json:"end_date,omitempty"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (r RecurrenceType) Valid() bool {
	switch r {
	case Daily, Monthly, SpecificDate, EvenDays, OddDays:
		return true
	default:
		return false
	}
}
