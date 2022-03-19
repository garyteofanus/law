package main

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo Repository
}

func (s service) getAllTask() ([]*Task, error) {
	return s.repo.GetTasks()
}

func (s service) getTask(id int) (*Task, error) {
	return s.repo.GetTaskByID(id)
}

func (s service) addTask(task *Task) (*Task, error) {
	if task.Title == "" {
		task.Title = "New Task"
	}

	createdTask, err := s.repo.CreateTask(task)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (s service) updateTask(id int, upd *Task) (*Task, error) {
	updatedTask, err := s.repo.UpdateTask(id, upd)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (s service) deleteTask(id int) error {
	err := s.repo.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   bool
}

type Service interface {
	getAllTask() ([]*Task, error)
	getTask(id int) (*Task, error)
	addTask(task *Task) (*Task, error)
	updateTask(id int, upd *Task) (*Task, error)
	deleteTask(id int) error
}
