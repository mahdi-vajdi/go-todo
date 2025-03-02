package todo

type Repository interface {
	Create(todo *Todo) error
	GetOneById(todoId int64, userId int64) (*Todo, error)
	UpdateCompleted(todoId int64, userId int64, completed bool) error
}
