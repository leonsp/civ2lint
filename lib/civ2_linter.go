package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func New(c Config, l *zap.SugaredLogger) Civ2Linter {
	cl := Civ2Linter{
		Config: c,
		Logger: l,
		Parser: RulesParser{
			Sections:     make(map[string][]LineInfo, 25),
			SeenSections: make(map[string]bool, 25),
		},
	}
	return cl
}

func (cl *Civ2Linter) Lint() error {
	var err error
	err = cl.parseFile("rules.txt")
	if err != nil {
		cl.Logger.Error("parsing failed:", zap.Error(err))
		return err
	}

	fmt.Println("Seen sections", cl.Parser.SeenSections)

	err = cl.LintAdvances()
	if err != nil {
		cl.Logger.Error("linting advances failed:", zap.Error(err))
		return err
	}

	fmt.Println(cl.Rules.Civilize)

	return nil
}

func (cl *Civ2Linter) parseFile(filename string) error {
	filePath := filepath.Join(cl.Config.Path, filename)
	_, err := os.Stat(filePath)
	if err != nil {
		cl.Logger.Error(fmt.Sprintf("%s does not exist:", filename), zap.Error(err))
		return err
	}

	readFile, err := os.Open(filePath)
	if err != nil {
		cl.Logger.Error(fmt.Sprintf("could not open %s:", filename), zap.Error(err))
		return err
	}
	defer func() { _ = readFile.Close() }()

	content, err := io.ReadAll(readFile)
	if err != nil {
		return err
	}

	parser := &RulesParser{}
	if err := parser.Parse(string(content)); err != nil {
		return err
	}

	cl.Parser = *parser
	return nil
}
