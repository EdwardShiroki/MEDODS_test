package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title           string                    `json:"title"`
	Description     string                    `json:"description"`
	Status          taskdomain.Status         `json:"status"`
	RecurrenceType  taskdomain.RecurrenceType `json:"recurrence_type,omitempty"`
	RecurrenceValue int                       `json:"recurrence_value,omitempty"`
	SpecificDates   []time.Time               `json:"specific_dates,omitempty"`
	EndDate         time.Time                 `json:"end_date,omitempty"`
}

type taskDTO struct {
	ID              int64                     `json:"id"`
	Title           string                    `json:"title"`
	Description     string                    `json:"description"`
	Status          taskdomain.Status         `json:"status"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	RecurrenceType  taskdomain.RecurrenceType `json:"recurrence_type,omitempty"`
	RecurrenceValue int                       `json:"recurrence_value,omitempty"`
	SpecificDates   []time.Time               `json:"specific_dates,omitempty"`
	EndDate         time.Time                 `json:"end_date,omitempty"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:              task.ID,
		Title:           task.Title,
		Description:     task.Description,
		Status:          task.Status,
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
		RecurrenceType:  task.RecurrenceType,
		RecurrenceValue: task.RecurrenceValue,
		SpecificDates:   task.SpecificDates,
		EndDate:         task.EndDate,
	}
}
