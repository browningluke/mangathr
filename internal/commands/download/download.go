package download

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"mangathrV2/internal/config"
	"mangathrV2/internal/sources/scrapers"
)

func SingleCheckboxes(label string, opts []string) string {
	var res string
	prompt := &survey.Select{
		Message: label,
		Options: opts,
	}
	err := survey.AskOne(prompt, &res, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = ""
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		panic(err)
	}

	return res
}

func Checkboxes(label string, opts []string) []string {
	res := []string{}
	prompt := &survey.MultiSelect{
		Message:  label,
		Options:  opts,
		PageSize: 10,
	}
	err := survey.AskOne(prompt, &res, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = ""
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		panic(err)
	}

	return res
}

func SelectManga(titles []string) string {
	selection := SingleCheckboxes(
		"Select Manga:",
		titles,
	)

	return selection
}

func SelectChapters(titles []string, mangaTitle string, sourceName string) []string {
	survey.MultiSelectQuestionTemplate = `
{{- define "option"}}
    {{- if eq .SelectedIndex .CurrentIndex }}{{color .Config.Icons.SelectFocus.Format }}{{ .Config.Icons.SelectFocus.Text }}{{color "reset"}}{{else}} {{end}}
    {{- if index .Checked .CurrentOpt.Index }}{{color .Config.Icons.MarkedOption.Format }} {{ .Config.Icons.MarkedOption.Text }} {{else}}{{color .Config.Icons.UnmarkedOption.Format }} {{ .Config.Icons.UnmarkedOption.Text }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{- .CurrentOpt.Value}}
{{end}}
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{"✔"}}{{color "reset"}}{{"\n"}}
{{- else }}
	{{- "  "}}{{- color "cyan"}}[Use arrows to move, space to select, <right> to all, <left> to none, type to filter{{- if and .Help (not .ShowHelp)}}, {{ .Config.HelpInput }} for more help{{end}}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
{{- end}}`

	selections := Checkboxes(
		fmt.Sprintf("\rTitle: %s\nSource: %s\n# of chapters: %d\nSelect chapters",
			mangaTitle, sourceName, len(titles)),
		titles,
	)

	return selections
}

func Run(args *Args, config *config.Config) {
	scraper := scrapers.NewScraper(args.Plugin)

	titles := scraper.Search(args.Query)
	//fmt.Println(titles)

	selection := SelectManga(titles)
	scraper.SelectManga(selection)

	_ = scraper.ListChapters()
	//chapterTitle := scraper.GetChapterTitle()
	//sourceName := scraper.GetScraperName()
	//chapterSelections := SelectChapters(chapters, chapterTitle, sourceName)
	//scraper.SelectChapters(chapterSelections)
}
