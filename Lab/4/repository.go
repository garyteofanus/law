package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sqlx.DB
}

func (r repository) GetTaskByID(id int) (*Task, error) {
	var task taskModel
	err := r.db.Get(&task, "SELECT * FROM tasks WHERE id = $1", id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to delete user with id %d: %w", id, err)
	}

	return &Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
	}, nil
}

func (r repository) GetTasks() ([]*Task, error) {
	var tasks []*taskModel
	query := "SELECT * FROM tasks"
	err := r.db.Select(&tasks, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	results := make([]*Task, len(tasks))
	for i, task := range tasks {
		results[i] = &Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		}
	}

	return results, nil
}

func (r repository) CreateTask(task *Task) (*Task, error) {
	query := "INSERT INTO tasks (title, description, completed) VALUES (:title, :description, :completed) RETURNING id"
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}

	var lastInsertID int
	err = stmt.Get(&lastInsertID, task)

	return &Task{
		ID:          lastInsertID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
	}, nil
}

func (r repository) UpdateTask(id int, upd *Task) (*Task, error) {
	query := "UPDATE tasks SET title=:title, description=:description, completed=:completed, updated_at=:updated_at WHERE id = :id"
	_, err := r.db.NamedExec(query, taskModel{
		ID:          id,
		Title:       upd.Title,
		Description: upd.Description,
		Completed:   upd.Completed,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return upd, nil
}

func (r repository) DeleteTask(id int) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return fmt.Errorf("task with id %d not found", id)
		}
		return fmt.Errorf("failed to delete task with id %d: %w", id, err)
	}

	return nil
}

type taskModel struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Completed   bool   `db:"completed"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Repository interface {
	GetTaskByID(id int) (*Task, error)
	GetTasks() ([]*Task, error)
	CreateTask(user *Task) (*Task, error)
	UpdateTask(id int, upd *Task) (*Task, error)
	DeleteTask(id int) error
}

