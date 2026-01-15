package task

import "time"

const (
	ToDo       = "todo"
	InProgress = "in-progress"
	Done       = "done"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreateAt    time.Time `json:"createAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewTask(Id int, Description, Status string) *Task {
	return &Task{
		Id:          Id,
		Description: Description,
		Status:      Status,
		CreateAt:    time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (t *Task) SetDescription(Description string) {
	t.Description = Description
}

func (t *Task) SetStatus(Status string) {
	t.Status = Status
}

func (t *Task) SetUpdatedAt(UpdatedAt time.Time) {
	t.UpdatedAt = UpdatedAt
}
