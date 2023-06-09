package homework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
)

const (
	apiSubjectsUrl        = "https://vshkole.com/api/get_class_subjects?new-app=1&class_id=%d&type=ab"
	apiSubjectEntitiesUrl = "https://vshkole.com/api/get_subject_class_entities?new-app=1&class_id=%d&subject_id=%s&type=ab"
	apiEntityUrl          = "https://vshkole.com/api/get_entity_by_id?new-app=1&id=%s&type=ab"
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

func (hg *homeworkGateway) GetSubjects(opts entity.Opts) ([]*entity.Subject, error) {
	var respData *GetSubjectsResponse
	if err := hg.makeRequest(context.Background(), fmt.Sprintf(apiSubjectsUrl, classToID[opts.Class]), &respData); err != nil {
		return nil, err
	}

	return respData.ToEntity(), nil
}

type SubjectEntity struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Author        string `json:"authors"`
	Year          string `json:"year"`
	Specification string `json:"specification"`
}

// TODO: make method cachable
func (hg *homeworkGateway) getSubjectEntities(opts entity.Opts) ([]*SubjectEntity, error) {
	var respData []*SubjectEntity
	err := hg.makeRequest(context.Background(), fmt.Sprintf(apiSubjectEntitiesUrl, classToID[opts.Class], opts.Subject.ID), &respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func (hg *homeworkGateway) GetAuthors(opts entity.Opts) ([]*entity.Author, error) {
	subjectEntities, err := hg.getSubjectEntities(opts)
	if err != nil {
		return nil, err
	}

	// TODO: unique authors: lower, strip
	var authors []*entity.Author
	for _, e := range subjectEntities {
		authors = append(authors, &entity.Author{Name: e.Author})
	}

	return authors, nil
}

func (hg *homeworkGateway) GetSpecifications(opts entity.Opts) ([]*entity.Specification, error) {
	subjectEntities, err := hg.getSubjectEntities(opts)
	if err != nil {
		return nil, err
	}

	var specifications []*entity.Specification
	for _, e := range subjectEntities {
		if e.Author != opts.Author.Name {
			continue
		}

		// TODO: maybe move to service layer?
		// TODO: it smells, refactor
		if e.Specification == "" {
			e.Specification = "Підручник"
		}
		// TODO: unique
		specifications = append(specifications, &entity.Specification{Name: e.Specification})
	}

	return specifications, nil
}

func (hg *homeworkGateway) GetYears(opts entity.Opts) ([]*entity.Year, error) {
	subjectEntities, err := hg.getSubjectEntities(opts)
	if err != nil {
		return nil, err
	}
	// TODO: it smells refactor
	if opts.Specification.Name == "Підручник" {
		opts.Specification.Name = ""
	}

	var years []*entity.Year
	for _, e := range subjectEntities {
		if e.Author != opts.Author.Name || e.Specification != opts.Specification.Name {
			continue
		}

		convertedYear, err := strconv.Atoi(e.Year)
		if err != nil {
			return nil, fmt.Errorf("convert year to int: %s", err)
		}
		// save subject entity ID to the last step that require subject entities (the "year" step)
		// TODO: maybe there's better way to save subject entity ID
		years = append(years, &entity.Year{ID: e.ID, Year: convertedYear})
	}

	return years, nil
}

type Entity struct {
	Contents []*EntityContent `json:"contents"`
}
type EntityContent struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Level     string           `json:"level"`
	Children  []*EntityContent `json:"__children"`
	ImageURLs []string         `json:"image_urls"`
	// Text string ?
}

// TODO: make method cachable
func (hg *homeworkGateway) getEntityContents(opts entity.Opts) ([]*EntityContent, error) {
	var respData *Entity
	err := hg.makeRequest(context.Background(), fmt.Sprintf(apiEntityUrl, opts.Year.ID), &respData)
	if err != nil {
		return nil, err
	}

	return respData.Contents, nil
}

func (hg *homeworkGateway) GetTopicsOrExercises(opts entity.Opts) ([]*entity.TopicOrExercise, error) {
	entityContents, err := hg.getEntityContents(opts)
	if err != nil {
		return nil, err
	}

	actualLevel := len(opts.Topics)
	for i := 0; i < actualLevel; i++ {
		for _, c := range entityContents {
			if c.ID != opts.Topics[i].ID {
				continue
			}
			entityContents = c.Children
		}
	}

	var res []*entity.TopicOrExercise
	for _, c := range entityContents {
		if len(c.Children) == 0 {
			res = append(res, &entity.TopicOrExercise{Exercise: &entity.Exercise{ID: c.ID, Name: c.Name}})
		} else {
			res = append(res, &entity.TopicOrExercise{Topic: &entity.Topic{ID: c.ID, Name: c.Name}})
		}
	}

	return res, nil
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
