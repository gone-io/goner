module github.com/gone-io/goner/mq/mqtt

go 1.24.1

require (
	dario.cat/mergo v1.0.2
	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/gone-io/gone/v2 v2.2.2
	github.com/stretchr/testify v1.9.0
	go.uber.org/mock v0.5.2
)

replace github.com/gone-io/goner/g => ../../g

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
