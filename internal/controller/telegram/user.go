package telegram

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"gopkg.in/telebot.v3"
)

var (
	InputSG            = fsm.NewStateGroup("reg")
	InputGrade         = InputSG.New("grade")
	InputSubject       = InputSG.New("subject")
	InputAuthor        = InputSG.New("author")
	InputSpecification = InputSG.New("spec")
	InputYear          = InputSG.New("year")
)

type HomeworkService interface {
	GetSubjects(opts entity.Opts) ([]*entity.Subject, error)
	GetAuthors(opts entity.Opts) ([]*entity.Author, error)
	GetSpecifications(opts entity.Opts) ([]*entity.Specification, error)
	GetYear(opts entity.Opts) ([]*entity.Year, error)
	GetTopics(opts entity.Opts) ([]*entity.Topic, error)
	GetExercises(opts entity.Opts) ([]*entity.Exercise, error)
}

type userHandler struct {
	homeworkService HomeworkService
}

func NewUserHandler(homeworkService HomeworkService) *userHandler {
	return &userHandler{homeworkService: homeworkService}
}

func (h *userHandler) Register(manager *fsm.Manager) {
	manager.Bind("/start", fsm.AnyState, h.OnStart)
	manager.Bind(telebot.OnText, InputGrade, h.OnInputGrade)
	manager.Bind(telebot.OnText, InputSubject, h.OnInputSubject)
	manager.Bind(telebot.OnText, InputAuthor, h.OnInputAuthor)
	manager.Bind(telebot.OnText, InputSpecification, h.OnInputSpecification)
	manager.Bind(telebot.OnText, InputYear, h.OnInputYear)
}

func (h *userHandler) OnStart(c telebot.Context, state fsm.Context) error {
	if err := state.Finish(true); err != nil {
		// TODO: handle
		return err
	}

	// TODO: create it not here
	markup := &telebot.ReplyMarkup{}
	markup.Reply(markup.Row(
		telebot.Btn{Text: "Back"}, // TODO
		telebot.Btn{Text: "1"},
		telebot.Btn{Text: "2"},
		telebot.Btn{Text: "3"},
		telebot.Btn{Text: "4"},
		telebot.Btn{Text: "5"},
		telebot.Btn{Text: "6"},
		telebot.Btn{Text: "7"},
		telebot.Btn{Text: "8"},
		telebot.Btn{Text: "9"},
		telebot.Btn{Text: "10"},
		telebot.Btn{Text: "11"},
	))

	if err := state.Set(InputGrade); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choice the grade", markup)
}

func (h *userHandler) OnInputGrade(c telebot.Context, state fsm.Context) error {
	log.Println("Grade:", c.Message().Text)

	grade, err := strconv.Atoi(c.Message().Text)
	if err != nil {
		return c.Send("You've put not a number, try again")
	}
	// TODO: input should be valid (number from 1 to 11)

	subjects, err := h.homeworkService.GetSubjects(entity.Opts{Grade: grade})
	if err != nil {
		return err
	}

	// TODO: create it not here
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, subject := range subjects {
		btns = append(btns, telebot.Btn{Text: subject.Name})
	}
	m.Reply(m.Split(4, btns)...)

	if err := state.Set(InputSubject); err != nil {
		// TODO: handle
		return err
	}
	if err := state.Update(InputGrade.String(), grade); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choose the subject", m)
}

func (h *userHandler) OnInputSubject(c telebot.Context, state fsm.Context) error {
	log.Println("Subject:", c.Message().Text)

	// TODO: input should be valid
	subject := &entity.Subject{Name: c.Message().Text}

	var grade int
	if err := getData(&getDataOpts{grade: &grade}, state); err != nil {
		// TODO: handle
		return err
	}

	authors, err := h.homeworkService.GetAuthors(entity.Opts{Grade: grade, Subject: subject})
	if err != nil {
		// TODO: handle
		return err

	}

	// TODO: create it not here
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, author := range authors {
		btns = append(btns, telebot.Btn{Text: author.Name})
	}
	m.Reply(m.Split(4, btns)...)

	if err := state.Update(InputSubject.String(), subject.Name); err != nil {
		// TODO: handle
		return err
	}
	if err := state.Set(InputAuthor); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choose the author", m)
}

func (h *userHandler) OnInputAuthor(c telebot.Context, state fsm.Context) error {
	log.Println("Author:", c.Message().Text)

	// TODO: input should be valid
	author := &entity.Author{Name: c.Message().Text}

	var (
		grade       int
		subjectName string
	)
	if err := getData(&getDataOpts{grade: &grade, subjectName: &subjectName}, state); err != nil {
		// TODO: handle
		return err
	}

	specifications, err := h.homeworkService.GetSpecifications(entity.Opts{
		Grade:   grade,
		Subject: &entity.Subject{Name: subjectName},
		Author:  author,
	})
	if err != nil {
		// TODO: handle
		return err
	}

	// TODO: create it not here
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, spec := range specifications {
		btns = append(btns, telebot.Btn{Text: spec.Name})
	}
	m.Reply(m.Split(4, btns)...)

	if err := state.Update(InputAuthor.String(), author.Name); err != nil {
		// TODO: handle
		return err
	}
	if err := state.Set(InputSpecification); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choose the specification", m)
}

func (h *userHandler) OnInputSpecification(c telebot.Context, state fsm.Context) error {
	log.Println("Specification:", c.Message().Text)

	// TODO: input should be valid
	specification := &entity.Specification{Name: c.Message().Text}

	var (
		grade       int
		subjectName string
		authorName  string
	)
	if err := getData(&getDataOpts{grade: &grade, subjectName: &subjectName, authorName: &authorName}, state); err != nil {
		// TODO: handle
		return err
	}

	years, err := h.homeworkService.GetYear(entity.Opts{
		Grade:         grade,
		Subject:       &entity.Subject{Name: subjectName},
		Author:        &entity.Author{Name: authorName},
		Specification: specification,
	})
	if err != nil {
		// TODO: handle
		return err
	}

	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, year := range years {
		btns = append(btns, telebot.Btn{Text: strconv.Itoa(year.Year)})
	}
	m.Reply(m.Split(4, btns)...)

	if err := state.Update(InputSpecification.String(), specification.Name); err != nil {
		// TODO: handle
		return err
	}
	if err := state.Set(InputYear); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choose the year", m)
}

func (h *userHandler) OnInputYear(c telebot.Context, state fsm.Context) error {
	log.Println("Year:", c.Message().Text)

	// TODO: input should be valid
	rawYear, err := strconv.Atoi(c.Message().Text)
	if err != nil {
		return c.Send("You've put not a number, try again")
	}
	year := &entity.Year{Year: rawYear}

	var (
		grade             int
		subjectName       string
		authorName        string
		specificationName string
	)
	if err := getData(&getDataOpts{grade: &grade, subjectName: &subjectName, authorName: &authorName, specificationName: &specificationName}, state); err != nil {
		// TODO: handle
		return err
	}
	log.Println("data:", grade, subjectName, authorName, specificationName)
	_ = year

	return nil
}

// TODO: better name
// TODO: move to other place
type getDataOpts struct {
	grade             *int
	subjectName       *string
	authorName        *string
	specificationName *string
}

// TODO: better name
// TODO: move to other place
// TODO: maybe we can avoid hardcoding the keys in `Get` calls to e.g. make it testable?
// gets the data and saves it to the passed pointers
func getData(opts *getDataOpts, state fsm.Context) error {
	if opts.grade != nil {
		if err := state.Get(InputGrade.String(), opts.grade); err != nil {
			return fmt.Errorf("get grade: %w", err)
		}
	}
	if opts.subjectName != nil {
		if err := state.Get(InputSubject.String(), opts.subjectName); err != nil {
			return fmt.Errorf("get subject name: %w", err)
		}
	}
	if opts.authorName != nil {
		if err := state.Get(InputAuthor.String(), opts.authorName); err != nil {
			return fmt.Errorf("get author name: %w", err)
		}
	}
	if opts.specificationName != nil {
		if err := state.Get(InputSpecification.String(), opts.specificationName); err != nil {
			return fmt.Errorf("get specification name: %w", err)
		}
	}
	return nil
}
