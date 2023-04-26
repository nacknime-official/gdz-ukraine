package service

import (
	"github.com/nacknime-official/gdz-ukraine/internal/entity"
	"golang.org/x/exp/slices"
)

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
	GetTopicsOrExercises(opts entity.Opts) ([]*entity.TopicOrExercise, error)
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
func (hs *homeworkService) GetSubjectByName(opts entity.Opts, name string) (*entity.Subject, error) {
	subjects, err := hs.GetSubjects(opts)
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(subjects, func(s *entity.Subject) bool { return s.Name == name })
	if idx == -1 {
		return nil, entity.ErrNotFound
	}

	return subjects[idx], nil
}

func (hs *homeworkService) GetAuthors(opts entity.Opts) ([]*entity.Author, error) {
	return hs.gateway.GetAuthors(opts)
}
func (hs *homeworkService) GetAuthorByName(opts entity.Opts, name string) (*entity.Author, error) {
	authors, err := hs.GetAuthors(opts)
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(authors, func(s *entity.Author) bool { return s.Name == name })
	if idx == -1 {
		return nil, entity.ErrNotFound
	}

	return authors[idx], nil
}

func (hs *homeworkService) GetSpecifications(opts entity.Opts) ([]*entity.Specification, error) {
	return hs.gateway.GetSpecifications(opts)
}
func (hs *homeworkService) GetSpecificationByName(opts entity.Opts, name string) (*entity.Specification, error) {
	specifications, err := hs.GetSpecifications(opts)
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(specifications, func(s *entity.Specification) bool { return s.Name == name })
	if idx == -1 {
		return nil, entity.ErrNotFound
	}

	return specifications[idx], nil
}

func (hs *homeworkService) GetYears(opts entity.Opts) ([]*entity.Year, error) {
	return hs.gateway.GetYears(opts)
}
func (hs *homeworkService) GetYearByValue(opts entity.Opts, year int) (*entity.Year, error) {
	years, err := hs.GetYears(opts)
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(years, func(s *entity.Year) bool { return s.Year == year })
	if idx == -1 {
		return nil, entity.ErrNotFound
	}

	return years[idx], nil
}

func (hs *homeworkService) GetTopicsOrExercises(opts entity.Opts) ([]*entity.TopicOrExercise, error) {
	return hs.gateway.GetTopicsOrExercises(opts)
}
func (hs *homeworkService) GetTopicOrExerciseByName(opts entity.Opts, name string) (*entity.TopicOrExercise, error) {
	topicsOrExercises, err := hs.GetTopicsOrExercises(opts)
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(topicsOrExercises, func(s *entity.TopicOrExercise) bool {
		if s.Topic != nil {
			return s.Topic.Name == name
		}
		if s.Exercise != nil {
			return s.Exercise.Name == name
		}
		return false
	})
	if idx == -1 {
		return nil, entity.ErrNotFound
	}

	return topicsOrExercises[idx], nil
}
