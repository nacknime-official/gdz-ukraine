package telegram

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"golang.org/x/exp/slices"
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
	GetYears(opts entity.Opts) ([]*entity.Year, error)
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

	var grade int
	if err := getData(&getDataOpts{grade: &grade}, state); err != nil {
		// TODO: handle
		return err
	}

	// check the input
	// TODO: refactor (check note 2 in "problems_notes.md")
	subjects, err := h.homeworkService.GetSubjects(entity.Opts{Grade: grade})
	if err != nil {
		// TODO: handle
		return err

	}
	idx := slices.IndexFunc(subjects, func(s *entity.Subject) bool { return s.Name == c.Message().Text })
	if idx < 0 {
		return c.Send("Click on one of the buttons!")
	}
	subject := subjects[idx]

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

	if err := state.Update(InputSubject.String(), subject); err != nil {
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
		grade   int
		subject entity.Subject
	)
	if err := getData(&getDataOpts{grade: &grade, subject: &subject}, state); err != nil {
		// TODO: handle
		return err
	}

	specifications, err := h.homeworkService.GetSpecifications(entity.Opts{
		Grade:   grade,
		Subject: &subject,
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

	if err := state.Update(InputAuthor.String(), author); err != nil {
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
		grade   int
		subject entity.Subject
		author  entity.Author
	)
	if err := getData(&getDataOpts{grade: &grade, subject: &subject, author: &author}, state); err != nil {
		// TODO: handle
		return err
	}

	years, err := h.homeworkService.GetYears(entity.Opts{
		Grade:         grade,
		Subject:       &subject,
		Author:        &author,
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

	if err := state.Update(InputSpecification.String(), specification); err != nil {
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
		grade         int
		subject       entity.Subject
		author        entity.Author
		specification entity.Specification
	)
	if err := getData(&getDataOpts{grade: &grade, subject: &subject, author: &author, specification: &specification}, state); err != nil {
		// TODO: handle
		return err
	}
	log.Println("data:", grade, subject, author, specification)
	_ = year

	return nil
}

// TODO: better name
// TODO: move to other place
type getDataOpts struct {
	grade         *int
	subject       *entity.Subject
	author        *entity.Author
	specification *entity.Specification
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
	if opts.subject != nil {
		if err := state.Get(InputSubject.String(), opts.subject); err != nil {
			return fmt.Errorf("get subject: %w", err)
		}
	}
	if opts.author != nil {
		if err := state.Get(InputAuthor.String(), opts.author); err != nil {
			return fmt.Errorf("get author: %w", err)
		}
	}
	if opts.specification != nil {
		if err := state.Get(InputSpecification.String(), opts.specification); err != nil {
			return fmt.Errorf("get specification: %w", err)
		}
	}
	return nil
}
