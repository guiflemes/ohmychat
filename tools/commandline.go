package tools

import (
	"bufio"
	"fmt"
	"oh-my-chat/src/service"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type myHandler struct{}

func (h *myHandler) Handle(message service.Message) error {
	fmt.Println("My handler")
	return nil
}

// USED TO TEST BUS AT COMMAND LINE
func RunBus() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	bus := service.NewBus()
	bus.SetHandler("some_name", &myHandler{})

	bus.Consume()

	commandLine(bus)

}

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
)

type SomeMessage struct{}

func (s *SomeMessage) Meta() service.MessageMeta {
	return service.MessageMeta{
		Id:    uuid.New(),
		Topic: "some_name",
	}
}

type myBus interface {
	Close()
	Publish(service.Message)
}

func commandLine(bus myBus) {

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

			bus.Publish(message)
		}

	}

}
