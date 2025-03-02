package todo

type Service interface {
	Create(todo *Todo) error
	GetById(todoId int64, userId int64) (*Todo, error)
	UpdateCompleted(todoId int64, userId int64, completed bool) error
}

type TodoService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &TodoService{repository: repository}
}

func (s *TodoService) Create(todo *Todo) error {
	err := s.repository.Create(todo)
	if err != nil {
		return err
	}

	return nil
}

func (s *TodoService) GetById(todoId int64, userId int64) (*Todo, error) {
	todo, err := s.repository.GetOneById(todoId, userId)
	if err != nil {
		return nil, err
	}

	return todo, err
}

func (s *TodoService) UpdateCompleted(todoId int64, userId int64, completed bool) error {
	err := s.repository.UpdateCompleted(todoId, userId, completed)
	if err != nil {
		return err
	}

	return nil
}
