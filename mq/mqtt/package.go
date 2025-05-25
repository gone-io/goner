package mqtt

//go:generate mockgen -destination=mqtt_mock.go -package=mqtt github.com/eclipse/paho.mqtt.golang Client,Token
