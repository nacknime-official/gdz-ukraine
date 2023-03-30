package homework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
)

const (
	apiSubjectsUrl        = "https://vshkole.com/api/get_class_subjects?new-app=1&class_id=%d&type=ab"
	apiSubjectEntitiesUrl = "https://vshkole.com/api/get_subject_class_entities?new-app=1&class_id=%d&subject_id=%d&type=ab"
	apiEntityUrl          = "https://vshkole.com/api/get_entity_by_id?new-app=1&id=%d&type=ab"
)

var classToID = map[int]int{
	1:  11,
	2:  10,
	3:  9,
	4:  8,
	5:  1,
	6:  2,
	7:  3,
	8:  4,
	9:  5,
	10: 6,
	11: 7,
}

type homeworkGateway struct {
	httpClient *http.Client
}

func NewHomeworkGateway(httpClient *http.Client) *homeworkGateway {
	return &homeworkGateway{httpClient: httpClient}
}

// TODO: pass contexts

func (hg *homeworkGateway) GetSubjects(opts entity.Opts) ([]*entity.Subject, error) {
	var respData *GetSubjectsResponse
	if err := hg.makeRequest(context.Background(), fmt.Sprintf(apiSubjectsUrl, classToID[opts.Class]), &respData); err != nil {
		return nil, err
	}

	return respData.ToEntity(), nil
}

type GetSubjectsResponse []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r GetSubjectsResponse) ToEntity() []*entity.Subject {
	var subjects []*entity.Subject
	for _, s := range r {
		subjects = append(subjects, &entity.Subject{ID: s.ID, Name: s.Name})
	}
	return subjects
}

func (hg *homeworkGateway) makeRequest(ctx context.Context, url string, decodeTo any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := hg.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&decodeTo); err != nil {
		return fmt.Errorf("json decode response: %w", err)
	}
	return nil
}

// TODO
func (hg *homeworkGateway) GetAuthors(opts entity.Opts) ([]*entity.Author, error) {
	return []*entity.Author{{"1", "1"}, {"2", "2"}, {"3", "3"}}, nil
}

// TODO
func (hg *homeworkGateway) GetSpecifications(opts entity.Opts) ([]*entity.Specification, error) {
	return []*entity.Specification{{"1", "Handbook"}, {"2", "Notebook"}}, nil
}

// TODO
func (hg *homeworkGateway) GetYears(opts entity.Opts) ([]*entity.Year, error) {
	return []*entity.Year{{"1", 2012}, {"2", 2015}, {"3", 2017}, {"4", 2022}}, nil
}

// TODO
func (hg *homeworkGateway) GetTopics(opts entity.Opts) ([]*entity.Topic, error) {
	return nil, nil
}

// TODO
func (hg *homeworkGateway) GetExercises(opts entity.Opts) ([]*entity.Exercise, error) {
	return nil, nil
}
