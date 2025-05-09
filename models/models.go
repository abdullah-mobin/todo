package models

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
type TodoList struct {
	Todos []Todo `json:"todos"`
}
