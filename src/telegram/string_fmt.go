package telegram

import (
	"fmt"
	"oh-my-chat/src/notion"
	"oh-my-chat/src/utils"
	"strings"
	"text/template"
)

type pendencyMessage struct {
	Name       string
	Greeting   string
	Pendencies []string
}

type pendencyFormatter struct{}

func (h *pendencyFormatter) Template() *template.Template {
	templateString := utils.NewStringBuilder().
		NextLine("Hello {{.Name}}, {{.Greeting}}").
		NextLine("These are your pending tasks:").
		NextLine("{{range .Pendencies}}").
		NextLine("- {{.}}").
		NextLine("{{end}}").
		String()

	tmpl, _ := template.New("message").Parse(templateString)
	return tmpl
}

func (h *pendencyFormatter) Format(pendency []notion.StudyStep) (string, error) {
	pendencies := make([]string, 0, len(pendency))
	for _, p := range pendency {
		text := fmt.Sprintf("[ %s ] %s - %s", p.Category, p.Name, p.CreatedAt.Format("02-01-2006"))
		pendencies = append(pendencies, text)
	}

	var result strings.Builder
	tmpl := h.Template()
	if err := tmpl.Execute(&result, pendencyMessage{Name: "Boss", Greeting: "morning", Pendencies: pendencies}); err != nil {
		fmt.Println("Error executing template:", err)
		return "", err
	}

	return result.String(), nil

}
