package main

import (
	"context"
	"fmt"
	"log"
	"storage/pkg/storage"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	DB_URL := "postgres://alexa:alexa@localhost:5432/tasks"

	db, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v\n", err)
	}
	defer db.Close()

	st := storage.NewStorage(db)

	fmt.Println("INSERT")
	t := storage.Task{
		AuthorID: 1,
		Title:    "First task",
		Content:  "какое-то текстовое описание",
		Labels:   []string{"one", "two"},
	}

	err = st.NewTask(&t)

	if err != nil {
		fmt.Println("failed:", err)
		return
	}
	fmt.Printf("task is: %#v\n", t)

	fmt.Println("UPDATE:")
	now := time.Now().Unix()
	t.Closed = &now
	err = st.UpdateTask(&t)
	fmt.Printf("task is: %#v\n", t)

	fmt.Println("SELECT ALL TASKS")
	ts, err := st.GetTasks()
	if err != nil {
		fmt.Println("failed: ", err)
		return
	}

	for _, t := range ts {
		fmt.Printf("%#v\n", t)
	}

	fmt.Println("GET TASKS BY AUTHOR")
	ts, err = st.GetTasksByAuthor(1)
	if err != nil {
		log.Fatalf("failed: %v\n", err)
	}

	for _, t := range ts {
		fmt.Printf("%#v\n", t)
	}

	if len(ts) == 0 {
		fmt.Println("задач такого автора нет")
	}

	fmt.Println("GET TASKS BY LABEL")
	ts, err = st.GetTasksByLabel("one")
	if err != nil {
		log.Fatalf("failed: %v\n", err)
	}

	for _, t := range ts {
		fmt.Printf("%#v\n", t)
	}

	fmt.Println("DELETE")
	err = st.DeleteTask(t.ID)
	if err != nil {
		fmt.Println("ошибка удаления", err)
	} else {
		fmt.Println("удалена")
	}

	err = st.DeleteTask(t.ID)
	if err != nil {
		fmt.Println("ошибка удаления:", err)
	}
}
