package lib

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func (cl *Civ2Linter) LintAdvances() error {
	section := "@CIVILIZE"
	lines, ok := cl.Parser.Sections[section]
	if !ok {
		message := fmt.Sprintf("section missing: %s", section)
		cl.Logger.Error(message)
		return errors.New(message)
	}

	cl.Rules.Civilize = make(map[string]Civilize, len(lines))
	for i, info := range lines {
		num := info.LineNumber
		line := info.Text
		comment_cols := strings.Split(line, ";")
		cols := strings.Split(comment_cols[0], ",")
		if len(cols) < 7 {
			return fmt.Errorf("too few columns in line %d: %s", num, line)
		} else if len(cols) > 7 {
			return fmt.Errorf("too many columns in line %d: %s", num, line)
		}
		aiValue, err := strconv.Atoi(strings.TrimSpace(cols[1]))
		if err != nil {
			return fmt.Errorf("invalid aiValue in line %d, %q: %w", num, line, err)
		}
		modifier, err := strconv.Atoi(strings.TrimSpace(cols[2]))
		if err != nil {
			return fmt.Errorf("invalid modifier in line %d, %q: %w", num, line, err)
		}
		epoch, err := strconv.Atoi(strings.TrimSpace(cols[5]))
		if err != nil {
			return fmt.Errorf("invalid epoch in line %d, %q: %w", num, line, err)
		}
		category, err := strconv.Atoi(strings.TrimSpace(cols[6]))
		if err != nil {
			return fmt.Errorf("invalid category in line %d, %q: %w", num, line, err)
		}
		advance := Civilize{
			Name:     strings.TrimSpace(cols[0]),
			AiValue:  aiValue,
			Modifier: modifier,
			Preq1:    strings.TrimSpace(cols[3]),
			Preq2:    strings.TrimSpace(cols[4]),
			Epoch:    epoch,
			Category: category,
		}
		cl.Rules.Civilize[AdvanceCodes[i]] = advance
	}

	for code, advance := range cl.Rules.Civilize {
		if code == advance.Preq1 || code == advance.Preq2 {
			cl.Rules.Errors = append(cl.Rules.Errors, fmt.Errorf("advance cannot be its own prerequisite: %s, %v", code, advance))
		}
		if (advance.Preq1 == "no" && advance.Preq2 != "no") || (advance.Preq1 != "no" && advance.Preq2 == "no") {
			cl.Rules.Errors = append(cl.Rules.Errors, fmt.Errorf("both prerequisites must be no: %s, %v", code, advance))
		}

		err := cl.FindLoops([]string{}, code)
		if err != nil {
			cl.Rules.Errors = append(cl.Rules.Errors, err)
		}
	}

	if len(cl.Rules.Errors) > 0 {
		return errors.Join(cl.Rules.Errors...)
	}
	return nil
}

func (cl *Civ2Linter) FindLoops(seen []string, next string) error {
	if next == "nil" || next == "no" || next == "..." {
		return nil
	}
	if slices.Contains(seen, next) {
		return fmt.Errorf("found loop: %v, %s", seen, next)
	}

	seen = append(seen, next)
	advance, ok := cl.Rules.Civilize[next]
	if !ok {
		return fmt.Errorf("advance does not exist: %s, %v", next, seen)
	}
	var err error
	err = cl.FindLoops(seen, advance.Preq1)
	if err != nil {
		return err
	}
	seen2 := make([]string, len(seen))
	copy(seen2, seen)
	err = cl.FindLoops(seen2, advance.Preq2)
	if err != nil {
		return err
	}
	return nil
}
