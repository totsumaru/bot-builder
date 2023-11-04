package main

import "github.com/totsumaru/bot-builder/app"

func main() {
	app.CreateMessageEvent([]string{"1", "2"}, []string{"1", "2"}, "キーワード1", "complete")
}
