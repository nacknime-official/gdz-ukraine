package telegram

import (
	"log"
	"strconv"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"gopkg.in/telebot.v3"
)

var (
	InputSG      = fsm.NewStateGroup("reg")
	InputGrade   = InputSG.New("grade")
	InputSubject = InputSG.New("subject")
	InputAuthor  = InputSG.New("author")
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
	log.Println("Grade: ", c.Message().Text)

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
	log.Println("Subject: ", c.Message().Text)

	// TODO: input should be valid
	subject := &entity.Subject{Name: c.Message().Text}

	var grade int
	if err := state.Get(InputGrade.String(), &grade); err != nil {
		// TODO: handle
		return err
	}
	log.Println(grade)

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
	return c.Send("Input the author", m)
}

func (h *userHandler) OnInputAuthor(c telebot.Context, state fsm.Context) error {
	log.Println("OnInputAuthor")
	return nil
}
