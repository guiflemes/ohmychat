package service

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

// USED TO TEST BUS AT COMMAND LINE
func RunBus() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	bus := NewBus()
	bus.SetHandler(Handler{Topic: "some_name", HandlerFunc: func(message Message) { fmt.Println("handler event") }})

	go bus.Consume()

	commandLine(bus)

}

type SomeMessage struct{}

func (s *SomeMessage) Meta() MessageMeta {
	return MessageMeta{
		Id:    uuid.New(),
		Topic: "some_name",
	}
}

func commandLine(bus *messageBus) {

	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	check := func(answer string) bool {
		return answer == "Y" || answer == "N"
	}

	reader := bufio.NewReader(os.Stdin)

	initialInput := func() (string, error) {

		answer, err := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		answer = strings.ToUpper(answer)

		if err != nil {
			return "", err
		}

		return answer, nil
	}

	var num int
	var err error

	defer func() {
		if err != nil {
			red.Printf("An error has ocurred %s\n", err)
			bus.Close()
			return
		}
		green.Println("Obrigado volte sempre!!")
	}()

	for {

		green.Print("Gostaria de publicar alguma mensagem? Y/N\n")
		answer, err := initialInput()

		if err != nil {
			break
		}

		if !check(answer) {
			yellow.Printf("Resposta '%s' invalida, tente novamente\n", answer)
			continue
		}

		if answer == "N" {
			break
		}

		green.Print("Quantos mensagen gostaria de publicar?")

		_, err = fmt.Scan(&num)

		if err != nil {
			break
		}

		for i := 0; i < num; i++ {
			message := &SomeMessage{}

			go bus.Publish(message)
		}

	}

}
