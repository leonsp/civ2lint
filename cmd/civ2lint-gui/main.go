package gui

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/leonsp/civ2lint/lib"
)

const appId = "com.github.leonsp.civ2lint"

func Main() {
	a := app.NewWithID(appId)
	w := a.NewWindow("civ2lint GUI")

	pathLabel := widget.NewLabel("No file selected")
	resultsArea := widget.NewMultiLineEntry()
	_ = resultsArea

	var selectedPath string

	browseBtn := widget.NewButton("Browse File", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			defer func() { _ = uri.Close() }()
			if uri == nil {
				return
			}

			selectedPath = ""
			displayPath := uri.URI().Path()
			if after, ok := strings.CutPrefix(displayPath, "file://"); ok {
				displayPath = after
			}
			selectedPath = displayPath
			pathLabel.SetText(displayPath)

		}, w)
	})

	lintBtn := widget.NewButton("Lint File", func() {
		if selectedPath == "" {
			dialog.ShowError(fmt.Errorf("no file selected"), w)
			return
		}

		content, err := os.ReadFile(selectedPath)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		parser := &lib.RulesParser{}
		if err := parser.Parse(string(content)); err != nil {
			dialog.ShowError(err, w)
			return
		}

		linter := lib.New(lib.Config{Path: "", Logger: nil}, nil)
		linter.Parser = *parser

		err = linter.LintAdvances()
		if err != nil {
			resultsArea.SetText(fmt.Sprintf("Errors found:\n%s", err.Error()))
		} else {
			resultsArea.SetText("No errors found!")
		}
	})

	content := container.NewVBox(
		browseBtn,
		pathLabel,
		lintBtn,
		widget.NewLabel("Results:"),
		resultsArea,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
