-- Add recurrence fields to tasks table
ALTER TABLE tasks ADD COLUMN recurrence_type VARCHAR(50);
ALTER TABLE tasks ADD COLUMN recurrence_value INTEGER;
ALTER TABLE tasks ADD COLUMN specific_dates TIMESTAMP[];
ALTER TABLE tasks ADD COLUMN end_date TIMESTAMP;