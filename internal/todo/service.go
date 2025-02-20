package todo

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(todo *Todo) error {
	err := s.repository.Create(todo)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetById(todoId int64, userId int64) (*Todo, error) {
	todo, err := s.repository.GetOneById(todoId, userId)
	if err != nil {
		return nil, err
	}

	return todo, err
}

func (s *Service) UpdateCompleted(todoId int64, userId int64, completed bool) error {
	err := s.repository.UpdateCompleted(todoId, userId, completed)
	if err != nil {
		return err
	}

	return nil
}
