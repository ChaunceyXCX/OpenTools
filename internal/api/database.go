package api

type DatabaseService struct{}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func (s *DatabaseService) Get(key string) (string, error) {
	return "", nil
}

func (s *DatabaseService) Put(key, value string) error {
	return nil
}

func (s *DatabaseService) Delete(key string) error {
	return nil
}
