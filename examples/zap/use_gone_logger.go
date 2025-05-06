package main

import "github.com/gone-io/gone/v2"

type UseGoneLogger struct {
	gone.Flag
	logger gone.Logger `gone:"*"`
}

func (u *UseGoneLogger) PrintLog() {
	u.logger.Infof("hello %s", "GONE IO")
}
