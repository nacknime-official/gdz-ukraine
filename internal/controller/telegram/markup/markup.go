package markup

import (
	"strconv"

	"github.com/nacknime-official/gdz-ukraine/internal/entity"
	"gopkg.in/telebot.v3"
)

func Classes() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for i := 1; i <= 11; i++ {
		btns = append(btns, telebot.Btn{Text: strconv.Itoa(i)})
	}
	m.Reply(m.Split(4, btns)...)
	return m
}

func Subjects(subjects []*entity.Subject) *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, subject := range subjects {
		btns = append(btns, telebot.Btn{Text: subject.Name})
	}
	m.Reply(m.Split(4, btns)...)
	return m
}

func Authors(authors []*entity.Author) *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, author := range authors {
		btns = append(btns, telebot.Btn{Text: author.Name})
	}
	m.Reply(m.Split(4, btns)...)
	return m
}

func Specifications(specifications []*entity.Specification) *telebot.ReplyMarkup {
	// TODO: create it not here
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, spec := range specifications {
		btns = append(btns, telebot.Btn{Text: spec.Name})
	}
	m.Reply(m.Split(4, btns)...)
	return m
}

func Years(years []*entity.Year) *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, year := range years {
		btns = append(btns, telebot.Btn{Text: strconv.Itoa(year.Year)})
	}
	m.Reply(m.Split(4, btns)...)
	return m
}

func TopicsOrExercises(topicsOrExercises []*entity.TopicOrExercise) *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{ResizeKeyboard: true}
	var btns []telebot.Btn
	for _, topicOrExercise := range topicsOrExercises {
		var text string
		if topicOrExercise.Topic != nil {
			text = topicOrExercise.Topic.Name
		}
		if topicOrExercise.Exercise != nil {
			text = topicOrExercise.Exercise.Name
		}
		btns = append(btns, telebot.Btn{Text: text})
	}
	m.Reply(m.Split(4, btns)...)
	return m
}
