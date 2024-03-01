package service

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

// USED TO TEST BUS AT COMMAND LINDE
func RunBus() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{}, 1)

	queue := Queue{message: make(chan Message, 0), done: done}

	bus := NewBus(queue)
	bus.SetEventHandler(EventHandler{Topic: "some_name", Handler: func() { fmt.Println("handler event") }})

	go bus.Handler()

	commandLine(queue)

}

func commandLine(queue Queue) {

	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

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
	var wg sync.WaitGroup

	for {

		green.Print("Gostaria de publicar alguma mensagem? Y/N\n")
		answer, err := initialInput()

		if err != nil {
			fmt.Println("error", err)
			break
		}

		if !check(answer) {
			red.Printf("Resposta '%s' invalida, tente novamente\n", answer)
			continue
		}

		if answer == "N" {
			break
		}

		green.Print("Quantos mensagen gostaria de publica?")

		_, err = fmt.Scan(&num)

		if err != nil {
			fmt.Println("error", err)
			break
		}

		for i := 0; i < num; i++ {
			message := &SomeMessage{}

			wg.Add(1)
			go func(msg Message) {
				queue.Publish(msg)
				wg.Done()
			}(message)
		}

	}

	go func() {
		wg.Wait()

	}()

	green.Println("Obrigado volte sempre")

}
