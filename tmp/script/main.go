package main

import (
	"encoding/json"
	"fmt"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain"
	"github.com/totsumaru/bot-builder/context/component/domain/button"
)

func main() {
	id, err := context.NewUUID()
	if err != nil {
		panic(err)
	}

	svID, err := context.NewDiscordID("1234567890")
	if err != nil {
		panic(err)
	}

	appID, err := context.NewUUID()
	if err != nil {
		panic(err)
	}

	core, err := domain.NewComponentCore(id, svID, appID)
	if err != nil {
		panic(err)
	}

	l, err := button.NewLabel("test")
	if err != nil {
		panic(err)
	}

	s, err := button.NewStyle(button.ButtonStylePrimary)
	if err != nil {
		panic(err)
	}

	b, err := button.NewButton(core, l, s, context.URL{})
	if err != nil {
		panic(err)
	}

	jsonb, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	bb := button.Button{}
	if err = json.Unmarshal(jsonb, &bb); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", string(jsonb))
}
