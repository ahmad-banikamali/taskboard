package domain

import "time"

type TaskID int64

type Task struct {
	id          TaskID
	teamID      TeamID
	title       Title
	description string
	status      Status
	assigneeID  *UserID
	dueDate     *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTask(team *Team, id TaskID, actor UserID, title string, description string, now time.Time) (*Task, error) {
	err := team.RequireCanEditTasks(actor)
	if err != nil {
		return nil, err
	}
	newTitle, err := NewTitle(title)
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
		assigneeID:  nil,
		status:      StatusTodo,
		updatedAt:   now.UTC(),
	}, nil
}

func (t *Task) ChangeStatus(team *Team, actor UserID, next Status, now time.Time) error {
	if err := t.checkHasEditPermission(team, actor); err != nil {
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

func (t *Task) ID() TaskID {
	return t.id
}

func (t *Task) Title() Title {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

func (t *Task) AssigneeID() (UserID, error) {
	if t.assigneeID == nil {
		return "", ErrTaskNotAssigned
	}
	return *t.assigneeID, nil
}

func (t *Task) DueDate() (time.Time, error) {
	if t.dueDate == nil {
		return time.Time{}, ErrTaskNotDueDate
	}
	return *t.dueDate, nil
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Task) Status() Status {
	return t.status
}

func (t *Task) Assign(team *Team, actor UserID, assignee UserID, now time.Time) error {
	if err := t.checkHasEditPermission(team, actor); err != nil {
		return err
	}
	if err := team.RequireMember(assignee); err != nil {
		return err
	}
	t.assigneeID = &assignee
	t.updatedAt = now.UTC()
	return nil
}

func (t *Task) Unassign(team *Team, actor UserID, now time.Time) error {
	if err := t.checkHasEditPermission(team, actor); err != nil {
		return err
	}
	t.assigneeID = nil
	t.updatedAt = now.UTC()
	return nil
}

func (t *Task) SetDueDate(team *Team, actor UserID, dueDate time.Time, now time.Time) error {
	if err := t.checkHasEditPermission(team, actor); err != nil {
		return err
	}
	dueDate = dueDate.UTC()
	t.dueDate = &dueDate
	t.updatedAt = now.UTC()
	return nil
}

func (t *Task) checkHasEditPermission(team *Team, actor UserID) error {
	if t.teamID != team.id {
		return ErrTaskNotInTeam
	}
	err := team.RequireCanEditTasks(actor)
	if err != nil {
		return err
	}
	return nil
}
