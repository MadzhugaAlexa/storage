package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID         int
	Opened     int64
	Closed     *int64
	AuthorID   int
	AssignedID *int
	Title      string
	Content    string
	Labels     []string
}

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) Storage {
	return Storage{
		db: db,
	}
}

// API пакета storage должен позволять:

// Создавать новые задачи
func (s *Storage) NewTask(t *Task) error {
	t.Opened = time.Now().Unix()

	row := s.db.QueryRow(
		context.Background(),
		"INSERT INTO tasks(author_id, title, content, opened) values($1, $2, $3, $4) returning id",
		t.AuthorID, t.Title, t.Content, t.Opened,
	)

	var id int
	err := row.Scan(&id)

	if err != nil {
		return err
	}

	t.ID = id

	if t.Labels == nil {
		return nil
	}

	for _, label := range t.Labels {
		// по условию: При этом для простоты таблицы пользователей и меток мы заполним самостоятельно.
		// значит метка уже существует в labels

		r := s.db.QueryRow(context.Background(), "select id from labels where name = $1", label)
		var labelID int
		err := r.Scan(&labelID)
		if err != nil {
			return err
		}

		_, err = s.db.Exec(context.Background(), "INSERT INTO tasks_labels(task_id, label_id) values($1, $2)",
			t.ID, labelID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) loadTasks(sql string, args ...any) ([]Task, error) {
	tasks := make([]Task, 0)

	rows, err := s.db.Query(context.Background(), sql, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		task := Task{}

		err = rows.Scan(
			&task.ID, &task.Opened, &task.Closed, &task.AuthorID, &task.AssignedID, &task.Title, &task.Content,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	rows, err = s.db.Query(context.Background(), "select id, name from labels")
	if err != nil {
		return nil, err
	}

	labels := make(map[int]string)
	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		labels[id] = name
	}

	rows, err = s.db.Query(context.Background(), "select task_id, label_id from tasks_labels")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var taskID int
		var labelID int

		err = rows.Scan(&taskID, &labelID)
		if err != nil {
			return nil, err
		}

		for index, t := range tasks {
			if t.ID == taskID {
				tasks[index].Labels = append(t.Labels, labels[labelID])
			}
		}
	}

	return tasks, nil
}

// Получать список всех задач,
func (s *Storage) GetTasks() ([]Task, error) {
	sql := "select id, opened, closed, author_id, assigned_id, title, content from tasks;"
	return s.loadTasks(sql)
}

// Получать список задач по автору,
func (s *Storage) GetTasksByAuthor(id int) ([]Task, error) {
	sql := "select id, opened, closed, author_id, assigned_id, title, content from tasks where author_id = $1"

	return s.loadTasks(sql, id)
}

// Получать список задач по метке,
func (s *Storage) GetTasksByLabel(label string) ([]Task, error) {
	sql := `
		select tasks.id, tasks.opened, tasks.closed, tasks.author_id, tasks.assigned_id, tasks.title, tasks.content 
		from tasks
		join tasks_labels on tasks_labels.task_id = tasks.id
		join labels on labels.id = tasks_labels.label_id
		where labels.name = $1;
	`

	return s.loadTasks(sql, label)
}

// Обновлять задачу по id,
func (s *Storage) UpdateTask(t *Task) error {
	sql := `UPDATE tasks set title = $1, 
		content = $2, opened = $3, closed = $4, 
		assignment_id = $5
		where id = $6
		;`

	_, err := s.db.Exec(
		context.Background(),
		sql,
		t.Title, t.Content, t.Opened, t.Closed, t.AssignedID, t.ID)

	if err != nil {
		return err
	}

	return nil
}

// Удалять задачу по id.
func (s *Storage) DeleteTask(id int) error {
	sql := "DELETE FROM tasks WHERE id = $1"

	tag, err := s.db.Exec(context.Background(), sql, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("не найдена запись")
	}

	return nil
}
