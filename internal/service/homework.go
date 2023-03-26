package service

import "github.com/nacknime-official/gdz-ukraine/internal/entity"

type mockHomeworkService struct{}

func NewMockHomeworkService() *mockHomeworkService {
	return &mockHomeworkService{}
}

func (hs *mockHomeworkService) GetSubjects(opts entity.Opts) ([]*entity.Subject, error) {
	return []*entity.Subject{{"1", "English"}, {"2", "Math"}, {"3", "PE"}, {"4", "Informatics"}, {"5", "History"}}, nil
}
func (hs *mockHomeworkService) GetAuthors(opts entity.Opts) ([]*entity.Author, error) {
	return []*entity.Author{{"1", "1"}, {"2", "2"}, {"3", "3"}}, nil
}
func (hs *mockHomeworkService) GetSpecifications(opts entity.Opts) ([]*entity.Specification, error) {
	return []*entity.Specification{{"1", "Handbook"}, {"2", "Notebook"}}, nil
}
func (hs *mockHomeworkService) GetYears(opts entity.Opts) ([]*entity.Year, error) {
	return []*entity.Year{{"1", 2012}, {"2", 2015}, {"3", 2017}, {"4", 2022}}, nil
}
func (hs *mockHomeworkService) GetTopics(opts entity.Opts) ([]*entity.Topic, error) {
	return nil, nil
}
func (hs *mockHomeworkService) GetExercises(opts entity.Opts) ([]*entity.Exercise, error) {
	return nil, nil
}

type HomeworkGateway interface {
	GetSubjects(opts entity.Opts) ([]*entity.Subject, error)
	GetAuthors(opts entity.Opts) ([]*entity.Author, error)
	GetSpecifications(opts entity.Opts) ([]*entity.Specification, error)
	GetYears(opts entity.Opts) ([]*entity.Year, error)
	GetTopics(opts entity.Opts) ([]*entity.Topic, error)
	GetExercises(opts entity.Opts) ([]*entity.Exercise, error)
}

type homeworkService struct {
	gateway HomeworkGateway
}

func NewHomeworkService(homeworkGateway HomeworkGateway) *homeworkService {
	return &homeworkService{gateway: homeworkGateway}
}

func (hs *homeworkService) GetSubjects(opts entity.Opts) ([]*entity.Subject, error) {
	return hs.gateway.GetSubjects(opts)
}
func (hs *homeworkService) GetAuthors(opts entity.Opts) ([]*entity.Author, error) {
	return []*entity.Author{{"1", "1"}, {"2", "2"}, {"3", "3"}}, nil
}
func (hs *homeworkService) GetSpecifications(opts entity.Opts) ([]*entity.Specification, error) {
	return []*entity.Specification{{"1", "Handbook"}, {"2", "Notebook"}}, nil
}
func (hs *homeworkService) GetYears(opts entity.Opts) ([]*entity.Year, error) {
	return []*entity.Year{{"1", 2012}, {"2", 2015}, {"3", 2017}, {"4", 2022}}, nil
}
func (hs *homeworkService) GetTopics(opts entity.Opts) ([]*entity.Topic, error) {
	return nil, nil
}
func (hs *homeworkService) GetExercises(opts entity.Opts) ([]*entity.Exercise, error) {
	return nil, nil
}
