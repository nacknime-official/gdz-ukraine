package telegram

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"gopkg.in/telebot.v3"
)

var (
	InputSG            = fsm.NewStateGroup("reg")
	InputClass         = InputSG.New("class")
	InputSubject       = InputSG.New("subject")
	InputAuthor        = InputSG.New("author")
	InputSpecification = InputSG.New("spec")
	InputYear          = InputSG.New("year")
	InputTopic         = InputSG.New("topic")
)

type HomeworkService interface {
	GetSubjects(opts entity.Opts) ([]*entity.Subject, error)
	GetSubjectByName(opts entity.Opts, name string) (*entity.Subject, error)

	GetAuthors(opts entity.Opts) ([]*entity.Author, error)
	GetAuthorByName(opts entity.Opts, name string) (*entity.Author, error)

	GetSpecifications(opts entity.Opts) ([]*entity.Specification, error)
	GetSpecificationByName(opts entity.Opts, name string) (*entity.Specification, error)

	GetYears(opts entity.Opts) ([]*entity.Year, error)
	GetYearByValue(opts entity.Opts, year int) (*entity.Year, error)

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
	manager.Bind(telebot.OnText, InputClass, h.OnInputClass)
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

	if err := state.Set(InputClass); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choice the class", markup)
}

func (h *userHandler) OnInputClass(c telebot.Context, state fsm.Context) error {
	log.Println("Class:", c.Message().Text)

	class, err := strconv.Atoi(c.Message().Text)
	if err != nil {
		return c.Send("You've put not a number, try again")
	}
	// TODO: input should be valid (number from 1 to 11)

	subjects, err := h.homeworkService.GetSubjects(entity.Opts{Class: class})
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
	if err := state.Update(InputClass.String(), class); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choose the subject", m)
}

func (h *userHandler) OnInputSubject(c telebot.Context, state fsm.Context) error {
	log.Println("Subject:", c.Message().Text)

	var class int
	if err := getData(&getDataOpts{class: &class}, state); err != nil {
		// TODO: handle
		return err
	}

	// check the input
	// TODO: refactor (check note 2 in "problems_notes.md")
	subject, err := h.homeworkService.GetSubjectByName(entity.Opts{Class: class}, c.Message().Text)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return c.Send("Click on one of the buttons!")
		}
		// TODO: handle
		return err
	}

	authors, err := h.homeworkService.GetAuthors(entity.Opts{Class: class, Subject: subject})
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

	var (
		class   int
		subject entity.Subject
	)
	if err := getData(&getDataOpts{class: &class, subject: &subject}, state); err != nil {
		// TODO: handle
		return err
	}
	opts := entity.Opts{Class: class, Subject: &subject}

	// check the input
	// TODO: refactor (check note 2 in "problems_notes.md")
	author, err := h.homeworkService.GetAuthorByName(opts, c.Message().Text)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return c.Send("Click on one of the buttons!")
		}
		// TODO: handle
		return err
	}
	opts.Author = author

	specifications, err := h.homeworkService.GetSpecifications(opts)
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

	var (
		class   int
		subject entity.Subject
		author  entity.Author
	)
	if err := getData(&getDataOpts{class: &class, subject: &subject, author: &author}, state); err != nil {
		// TODO: handle
		return err
	}
	opts := entity.Opts{Class: class, Subject: &subject, Author: &author}

	// check the input
	// TODO: refactor (check note 2 in "problems_notes.md")
	specification, err := h.homeworkService.GetSpecificationByName(opts, c.Message().Text)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return c.Send("Click on one of the buttons!")
		}
		// TODO: handle
		return err
	}
	opts.Specification = specification

	years, err := h.homeworkService.GetYears(opts)
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

	rawYear, err := strconv.Atoi(c.Message().Text)
	if err != nil {
		return c.Send("You've put not a number, try again")
	}

	var (
		class         int
		subject       entity.Subject
		author        entity.Author
		specification entity.Specification
	)
	if err := getData(&getDataOpts{class: &class, subject: &subject, author: &author, specification: &specification}, state); err != nil {
		// TODO: handle
		return err
	}
	opts := entity.Opts{Class: class, Subject: &subject, Author: &author, Specification: &specification}

	year, err := h.homeworkService.GetYearByValue(opts, rawYear)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return c.Send("Click on one of the buttons!")
		}
		// TODO: handle
		return err
	}
	opts.Year = year

	topics, err := h.homeworkService.GetTopics(opts)
	if err != nil {
		// TODO: handle
		return err
	}

	// TODO: create it not here
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, topic := range topics {
		btns = append(btns, telebot.Btn{Text: topic.Name})
	}
	m.Reply(m.Split(4, btns)...)

	if err := state.Update(InputYear.String(), year); err != nil {
		// TODO: handle
		return err
	}
	if err := state.Set(InputTopic); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choose the topic", m)
}

// TODO: better name
// TODO: move to other place
type getDataOpts struct {
	class         *int
	subject       *entity.Subject
	author        *entity.Author
	specification *entity.Specification
}

// TODO: better name
// TODO: move to other place
// TODO: maybe we can avoid hardcoding the keys in `Get` calls to e.g. make it testable?
// gets the data and saves it to the passed pointers
func getData(opts *getDataOpts, state fsm.Context) error {
	if opts.class != nil {
		if err := state.Get(InputClass.String(), opts.class); err != nil {
			return fmt.Errorf("get class: %w", err)
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
