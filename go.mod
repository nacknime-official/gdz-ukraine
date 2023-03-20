module github.com/nacknime-official/gdz-ukraine

go 1.20

require (
	github.com/redis/go-redis/v9 v9.0.2
	github.com/vitaliy-ukiru/fsm-telebot v0.3.1
	github.com/vitaliy-ukiru/fsm-telebot/storages/redis v0.3.1
	golang.org/x/exp v0.0.0-20230315142452-642cacee5cc0
	gopkg.in/telebot.v3 v3.1.2
)

replace github.com/vitaliy-ukiru/fsm-telebot => ./fsm-telebot

replace github.com/vitaliy-ukiru/fsm-telebot/storages/redis => ./fsm-telebot/storages/redis

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
