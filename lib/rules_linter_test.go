package lib_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/leonsp/civ2lint/lib"
)

var sugar *zap.SugaredLogger

var _ = Describe("RuleLinter", func() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }() // flushes buffer, if any
	sugar = logger.Sugar()

	Describe("Linting advances", func() {
		It("detects missing section", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("section missing"))
		})

		It("detects missing columns", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {{Text: "Advanced Flight,    4,"}},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("too few columns in line"))
		})

		It("detects extra columns", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {{Text: "Advanced Flight,    4,-2,  nil, no, 3, 4, 9000"}},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("too many columns in line"))
		})

		It("detects invalid ai value", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {{Text: "Advanced Flight,    Shrek,-2,  nil, no, 3, 4"}},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid aiValue in line"))
		})

		It("detects invalid ai value", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {{Text: "Advanced Flight,    -4,Shrek,  nil, no, 3, 4"}},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid modifier in line"))
		})

		It("detects invalid ai value", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {{Text: "Advanced Flight,    -4,-2,  nil, no, Shrek, 4"}},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid epoch in line"))
		})

		It("detects invalid ai value", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {{Text: "Advanced Flight,    -4,-2,  nil, no, 3, Shrek"}},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid category in line"))
		})

		It("Detects partial no", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {{
							LineNumber: 20,
							Text:       "Advanced Flight,    4,-2,  nil, no, 3, 4",
							Comment:    "Hello world",
						}},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("both prerequisites must be no"))
		})
		It("Detects self-reference", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {
							{Text: "Advanced Flight,    4,-2,  AFl, nil, 3, 4"},
						},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("advance cannot be its own prerequisite"))
		})

		It("Detects loops on the first prerequisite", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {
							{Text: "Advanced Flight,    4,-2,  Alp, nil, 3, 4"},
							{Text: "Alphabet,           5, 1,  AFl, nil, 0, 3"},
						},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("found loop"))
		})

		It("detects loops on the second prerequisite", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {
							{Text: "Advanced Flight,    4,-2,  nil, Alp, 3, 4"},
							{Text: "Alphabet,           5, 1,  nil, AFl, 0, 3"},
						},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("found loop"))
		})

		It("allows referencing later advances", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {
							{Text: "Advanced Flight,    4,-2,  Alp, nil, 3, 4"},
							{Text: "Alphabet,           5, 1,  nil, nil, 0, 3"},
						},
					},
				},
			}

			Expect(cl.LintAdvances()).NotTo(HaveOccurred())
		})

		It("should include line content in error messages for loops", func() {
			cl := lib.Civ2Linter{
				Logger: sugar,
				Parser: lib.RulesParser{
					Sections: map[string][]lib.LineInfo{
						"@CIVILIZE": {
							{Text: "Extra Advance 6,3,0,X7,no,0,0; X6"},
							{Text: "Extra Advance 7,3,0,X6,no,0,0; X7"},
						},
					},
				},
			}

			err := cl.LintAdvances()
			Expect(err).To(HaveOccurred())
			// We expect the error message to contain the line content as defined in the spec
			Expect(err.Error()).To(ContainSubstring("Extra Advance 6 3 0 X7 no 0 0"))
		})
	})
})
