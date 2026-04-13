package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:           normalized.Title,
		Description:     normalized.Description,
		Status:          normalized.Status,
		RecurrenceType:  normalized.RecurrenceType,
		RecurrenceValue: normalized.RecurrenceValue,
		SpecificDates:   normalized.SpecificDates,
		EndDate:         normalized.EndDate,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	// If recurring, generate multiple tasks
	if model.RecurrenceType != "" {
		tasks, err := s.generateRecurringTasks(model, now)
		if err != nil {
			return nil, err
		}
		for _, t := range tasks {
			_, err := s.repo.Create(ctx, &t)
			if err != nil {
				return nil, err
			}
		}
		// Return the first created task as representative
		return &tasks[0], nil
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:              id,
		Title:           normalized.Title,
		Description:     normalized.Description,
		Status:          normalized.Status,
		RecurrenceType:  normalized.RecurrenceType,
		RecurrenceValue: normalized.RecurrenceValue,
		SpecificDates:   normalized.SpecificDates,
		EndDate:         normalized.EndDate,
		UpdatedAt:       s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceType != "" && !input.RecurrenceType.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
	}

	// Validate recurrence fields
	if input.RecurrenceType == taskdomain.Daily && input.RecurrenceValue <= 0 {
		return CreateInput{}, fmt.Errorf("%w: recurrence value must be positive for daily", ErrInvalidInput)
	}
	if input.RecurrenceType == taskdomain.Monthly && (input.RecurrenceValue < 1 || input.RecurrenceValue > 30) {
		return CreateInput{}, fmt.Errorf("%w: recurrence value must be 1-30 for monthly", ErrInvalidInput)
	}
	if input.RecurrenceType == taskdomain.SpecificDate && len(input.SpecificDates) == 0 {
		return CreateInput{}, fmt.Errorf("%w: specific dates required", ErrInvalidInput)
	}
	if input.EndDate.IsZero() == false && input.EndDate.Before(time.Now()) {
		return CreateInput{}, fmt.Errorf("%w: end date must be in future", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceType != "" && !input.RecurrenceType.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
	}

	if input.RecurrenceType == taskdomain.Daily && input.RecurrenceValue <= 0 {
		return UpdateInput{}, fmt.Errorf("%w: recurrence value must be positive for daily", ErrInvalidInput)
	}
	if input.RecurrenceType == taskdomain.Monthly && (input.RecurrenceValue < 1 || input.RecurrenceValue > 30) {
		return UpdateInput{}, fmt.Errorf("%w: recurrence value must be 1-30 for monthly", ErrInvalidInput)
	}
	if input.RecurrenceType == taskdomain.SpecificDate && len(input.SpecificDates) == 0 {
		return UpdateInput{}, fmt.Errorf("%w: specific dates required", ErrInvalidInput)
	}
	if input.EndDate.IsZero() == false && input.EndDate.Before(time.Now()) {
		return UpdateInput{}, fmt.Errorf("%w: end date must be in future", ErrInvalidInput)
	}

	return input, nil
}

func (s *Service) generateRecurringTasks(template *taskdomain.Task, now time.Time) ([]taskdomain.Task, error) {
	var tasks []taskdomain.Task
	end := template.EndDate
	if end.IsZero() {
		end = now.AddDate(0, 1, 0)
	}

	switch template.RecurrenceType {
	case taskdomain.Daily:
		for d := now; d.Before(end); d = d.AddDate(0, 0, template.RecurrenceValue) {
			task := *template
			task.CreatedAt = d
			task.UpdatedAt = d
			tasks = append(tasks, task)
		}
	case taskdomain.Monthly:
		for d := now; d.Before(end); d = d.AddDate(0, 1, 0) {
			taskDate := time.Date(d.Year(), d.Month(), template.RecurrenceValue, 0, 0, 0, 0, d.Location())
			if taskDate.After(now) && taskDate.Before(end) {
				task := *template
				task.CreatedAt = taskDate
				task.UpdatedAt = taskDate
				tasks = append(tasks, task)
			}
		}
	case taskdomain.SpecificDate:
		for _, date := range template.SpecificDates {
			if date.After(now) && date.Before(end) {
				task := *template
				task.CreatedAt = date
				task.UpdatedAt = date
				tasks = append(tasks, task)
			}
		}
	case taskdomain.EvenDays:
		for d := now; d.Before(end); d = d.AddDate(0, 0, 1) {
			if d.Day()%2 == 0 {
				task := *template
				task.CreatedAt = d
				task.UpdatedAt = d
				tasks = append(tasks, task)
			}
		}
	case taskdomain.OddDays:
		for d := now; d.Before(end); d = d.AddDate(0, 0, 1) {
			if d.Day()%2 != 0 {
				task := *template
				task.CreatedAt = d
				task.UpdatedAt = d
				tasks = append(tasks, task)
			}
		}
	}
	return tasks, nil
}
