package service

import "github.com/nacknime-official/gdz-ukraine/internal/entity"

type mockHomeworkService struct{}

func NewMockHomeworkService() *mockHomeworkService {
	return &mockHomeworkService{}
}

func (hs *mockHomeworkService) GetSubjects(opts entity.Opts) ([]*entity.Subject, error) {
	return []*entity.Subject{{"English"}, {"Math"}, {"PE"}, {"Informatics"}, {"History"}}, nil
}
func (hs *mockHomeworkService) GetAuthors(opts entity.Opts) ([]*entity.Author, error) {
	return []*entity.Author{{"1"}, {"2"}, {"3"}}, nil
}
func (hs *mockHomeworkService) GetSpecifications(opts entity.Opts) ([]*entity.Specification, error) {
	return nil, nil
}
func (hs *mockHomeworkService) GetYear(opts entity.Opts) ([]*entity.Year, error) {
	return nil, nil
}
func (hs *mockHomeworkService) GetTopics(opts entity.Opts) ([]*entity.Topic, error) {
	return nil, nil
}
func (hs *mockHomeworkService) GetExercises(opts entity.Opts) ([]*entity.Exercise, error) {
	return nil, nil
}
