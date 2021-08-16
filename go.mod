module nknovh

go 1.15

replace (
	nknovh-engine v1.1.0 => ./internal/nknovh-engine
	nknovh-wasm v1.0.0 => ./internal/nknovh-wasm
	xwasmapi v1.0.0 => ./internal/xwasmapi
	templater v1.0.0 => ./internal/templater
)

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/nknorg/nkn-sdk-go v1.3.6 // indirect
	github.com/sevlyar/go-daemon v0.1.5 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	nknovh-engine v1.1.0 // indirect
	nknovh-wasm v1.0.0
	xwasmapi v1.0.0
	templater v1.0.0
)
