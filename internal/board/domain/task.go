package domain

import "time"

type TaskId int64

type Task struct {
	id          TaskId
	teamID      TeamID
	title       Title
	description string
	status      Status
	assigneeID  *UserID
	dueDate     *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTask(team *Team, id TaskId, actor UserID, title string, description string, now time.Time) (*Task, error) {
	newTitle, err := NewTitle(title)
	if err != nil {
		return nil, err
	}
	err = team.RequireCanEditTasks(actor)
	if err != nil {
		return nil, err
	}
	return &Task{
		id:          id,
		teamID:      team.id,
		title:       newTitle,
		description: description,
		createdAt:   now.UTC(),
		dueDate:     nil,
		assigneeID:  &actor,
		status:      StatusTodo,
	}, nil
}

func (t *Task) ChangeStatus(team *Team, actor UserID, next Status, now time.Time) error {
	err := team.requireOwner(actor)
	if err != nil {
		return err
	}

	status, err := t.status.MoveTo(next)
	if err != nil {
		return err
	}

	t.status = status
	t.updatedAt = now.UTC()
	return nil
}
