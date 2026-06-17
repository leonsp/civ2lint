package lib

import (
	"bufio"
	"strings"
)

func (p *RulesParser) Parse(content string) error {
	p.Sections = make(map[string][]LineInfo)
	p.SeenSections = make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentSection string
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		text, comment, _ := strings.Cut(line, ";")
		trimmed := strings.TrimSpace(text)

		if strings.HasPrefix(trimmed, "@") {
			currentSection = trimmed
			p.Sections[currentSection] = []LineInfo{}
		} else if currentSection != "" && len(trimmed) > 0 {
			p.Sections[currentSection] = append(p.Sections[currentSection], LineInfo{
				LineNumber: lineNum,
				Text:       line,
				Comment:    comment,
			})
		}
	}

	for section := range p.Sections {
		p.SeenSections[section] = true
	}
	return scanner.Err()
}
