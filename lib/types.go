package lib

import (
	"go.uber.org/zap"
)

type Config struct {
	Path   string
	Logger *zap.SugaredLogger
}

type Civ2Linter struct {
	Config Config
	Logger *zap.SugaredLogger
	Parser RulesParser
	Rules  Civ2Rules
}

type Civ2Rules struct {
	Civilize map[string]Civilize
	Errors   []error
}

type Civilize struct {
	Name     string
	AiValue  int
	Modifier int
	Preq1    string
	Preq2    string
	Epoch    int
	Category int
}

type LineInfo struct {
	LineNumber int
	Text       string
	Comment    string
}

type RulesParser struct {
	Sections     map[string][]LineInfo
	SeenSections map[string]bool
}
