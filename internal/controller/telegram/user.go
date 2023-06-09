package telegram

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/nacknime-official/gdz-ukraine/internal/controller/telegram/markup"
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
	InputTopics        = InputSG.New("topics")
	InputExercise      = InputSG.New("exercise")
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

	GetTopicsOrExercises(opts entity.Opts) ([]*entity.TopicOrExercise, error)
	GetTopicOrExerciseByName(opts entity.Opts, name string) (*entity.TopicOrExercise, error)
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
	manager.Bind(telebot.OnText, InputTopics, h.OnInputTopic)
}

func (h *userHandler) OnStart(c telebot.Context, state fsm.Context) error {
	if err := state.Finish(true); err != nil {
		// TODO: handle
		return err
	}

	if err := state.Set(InputClass); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choice the class", markup.Classes())
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
	m := markup.Subjects(subjects)

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
	m := markup.Authors(authors)

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
	m := markup.Specifications(specifications)

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
	m := markup.Years(years)

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

	topicsOrExercises, err := h.homeworkService.GetTopicsOrExercises(opts)
	if err != nil {
		// TODO: handle
		return err
	}

	m := markup.TopicsOrExercises(topicsOrExercises)

	if err := state.Update(InputYear.String(), year); err != nil {
		// TODO: handle
		return err
	}
	if err := state.Set(InputTopics); err != nil {
		// TODO: handle
		return err
	}
	return c.Send("Choose the topic", m)
}

func (h *userHandler) OnInputTopic(c telebot.Context, state fsm.Context) error {
	log.Println("Topic:", c.Message().Text)

	var (
		class         int
		subject       entity.Subject
		author        entity.Author
		specification entity.Specification
		year          entity.Year
		topics        []*entity.Topic
	)
	if err := getData(&getDataOpts{class: &class, subject: &subject, author: &author, specification: &specification, year: &year}, state); err != nil {
		// TODO: handle
		return err
	}
	// handle case when topic has not set yet
	if err := getData(&getDataOpts{topics: &topics}, state); err != nil {
		if !errors.Is(err, fsm.ErrNotFound) {
			// TODO: handle
			return err
		}
	}
	opts := entity.Opts{Class: class, Subject: &subject, Author: &author, Specification: &specification, Year: &year, Topics: topics}

	// check the input
	topicOrExercise, err := h.homeworkService.GetTopicOrExerciseByName(opts, c.Message().Text)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return c.Send("Click on one of the buttons!")
		}
		// TODO: handle
		return err
	}
	if topicOrExercise.Exercise != nil {
		// it's exercise, so we need to get the solution
		// but firstly save the exercise to the state
		// but don't sure if it's needed
		if err := state.Update(InputExercise.String(), topicOrExercise.Exercise); err != nil {
			// TODO: handle
			return err
		}
		// TODO: send solution
		return nil
	}
	// it's topic, so we need to get the next topics
	opts.Topics = append(opts.Topics, topicOrExercise.Topic)

	nextTopics, err := h.homeworkService.GetTopicsOrExercises(opts)
	if err != nil {
		// TODO: handle
		return err
	}
	if len(nextTopics) != 0 {
		m := markup.TopicsOrExercises(nextTopics)
		if err := state.Update(InputTopics.String(), opts.Topics); err != nil {
			// TODO: handle
			return err
		}
		return c.Send("Choose the topic", m)
	} else {
		return c.Send("weird")
	}
}

// TODO: better name
// TODO: move to other place
type getDataOpts struct {
	class         *int
	subject       *entity.Subject
	author        *entity.Author
	specification *entity.Specification
	year          *entity.Year
	topics        *[]*entity.Topic
	exercise      *entity.Exercise
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
	if opts.year != nil {
		if err := state.Get(InputYear.String(), opts.year); err != nil {
			return fmt.Errorf("get year: %w", err)
		}
	}
	if opts.topics != nil {
		if err := state.Get(InputTopics.String(), opts.topics); err != nil {
			return fmt.Errorf("get topics: %w", err)
		}
	}
	if opts.exercise != nil {
		if err := state.Get(InputExercise.String(), opts.exercise); err != nil {
			return fmt.Errorf("get exercise: %w", err)
		}
	}
	return nil
}
