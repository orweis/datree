package bl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/printer"
	"github.com/datreeio/datree/pkg/propertiesExtractor"
	"gopkg.in/yaml.v3"
)

type Printer interface {
	PrintWarnings(warnings []printer.Warning)
	PrintSummaryTable(summary printer.Summary)
}
type CLIClient interface {
	RequestEvaluation(pattern string, files []*propertiesExtractor.FileProperties, cliId string) (cliClient.EvaluationResponse, error)
}
type PropertiesExtractor interface {
	ReadFilesFromPattern(pattern string, conc int) ([]*propertiesExtractor.FileProperties, []propertiesExtractor.FileError, []error)
}
type Evaluator struct {
	propertiesExtractor PropertiesExtractor
	cliClient           CLIClient
	printer             Printer
}

func CreateNewEvaluator(pe PropertiesExtractor, c CLIClient, p Printer) *Evaluator {
	return &Evaluator{
		propertiesExtractor: pe,
		cliClient:           c,
		printer:             p,
	}
}

type EvaluationResults struct {
	FileNameRuleMapper map[string]map[int]*Rule
	Summary            struct {
		RulesCount       int
		TotalFailedRules int
		FilesCount       int
	}
}

func (e *Evaluator) Evaluate(pattern string, cliId string, evaluationConc int) (*EvaluationResults, []propertiesExtractor.FileError, error) {
	files, fileErrors, errors := e.propertiesExtractor.ReadFilesFromPattern(pattern, evaluationConc)
	if len(errors) > 0 {
		return nil, fileErrors, fmt.Errorf("failed evaluation with the following errors: %s", errors)
	}

	if len(files) == 0 {
		return nil, fileErrors, fmt.Errorf("no files detected")
	}

	res, err := e.cliClient.RequestEvaluation(pattern, files, cliId)
	if err != nil {
		return nil, fileErrors, err
	}

	results := e.aggregateEvaluationResults(res.Results, len(files))

	return results, fileErrors, nil
}

func (e *Evaluator) PrintResults(results *EvaluationResults, cliId string, output string) error {
	if output == "json" {
		jsonOutput, err := json.Marshal(results)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(string(jsonOutput))
		return nil
	}

	if output == "yaml" {
		yamlOutput, err := yaml.Marshal(results)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(string(yamlOutput))
		return nil
	}

	warnings, err := e.parseEvaluationResultsToWarnings(results)
	if err != nil {
		fmt.Println(err)
		return err
	}
	e.printer.PrintWarnings(warnings)

	configurePolicyLink := "https://app.datree.io/login?cliId=" + cliId
	summary := e.parseEvaluationResultsToSummary(results, configurePolicyLink)

	e.printer.PrintSummaryTable(summary)

	if results.Summary.TotalFailedRules > 0 {
		return fmt.Errorf("failed rules count is %d (>0)", results.Summary.TotalFailedRules)
	}
	return nil
}

func (e *Evaluator) PrintFileParsingErrors(errors []propertiesExtractor.FileError) {
	if len(errors) > 0 {
		fmt.Println("The following files failed:")

		for _, fileError := range errors {
			fmt.Printf("\n\tFilename: %s\n\tError: %s\n\t---------------------", fileError.Filename, fileError.Message)
		}
	}
}

func (e *Evaluator) aggregateEvaluationResults(evaluationResults []cliClient.EvaluationResult, filesCount int) *EvaluationResults {
	mapper := make(map[string]map[int]*Rule)

	totalRulesCount := len(evaluationResults)
	totalFailedCount := 0
	filenames := make(map[string]string)

	for _, result := range evaluationResults {
		if !result.Passed {
			totalFailedCount++
		}
		for _, match := range result.Results.Matches {
			// file not already exists in mapper
			if _, exists := mapper[match.FileName]; !exists {
				mapper[match.FileName] = make(map[int]*Rule)
			}

			// file and rule not already exists in mapper
			if _, exists := mapper[match.FileName][result.Rule.ID]; !exists {
				filenames[match.FileName] = match.FileName
				mapper[match.FileName][result.Rule.ID] = &Rule{ID: result.Rule.ID, Name: result.Rule.Name, FailSuggestion: result.Rule.FailSuggestion, Count: 0}
			}

			mapper[match.FileName][result.Rule.ID].IncrementCount()
		}
		for _, mismatch := range result.Results.Mismatches {
			if _, exists := mapper[mismatch.FileName]; !exists {
				filenames[mismatch.FileName] = mismatch.FileName
			}
		}
	}

	results := &EvaluationResults{
		FileNameRuleMapper: mapper,
		Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
		}{
			RulesCount:       totalRulesCount,
			TotalFailedRules: totalFailedCount,
			FilesCount:       filesCount,
		},
	}

	return results
}

func (e *Evaluator) parseEvaluationResultsToWarnings(results *EvaluationResults) ([]printer.Warning, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	warnings := []printer.Warning{}

	for fileName, rules := range results.FileNameRuleMapper {
		relativePath, _ := filepath.Rel(pwd, fileName)
		title := fmt.Sprintf(">>  File: %s\n", relativePath)

		warning := printer.Warning{
			Title: title,
		}

		for _, rule := range rules {
			details := struct {
				Caption     string
				Occurrences int
				Suggestion  string
			}{
				Caption:     rule.Name,
				Occurrences: rule.Count,
				Suggestion:  rule.FailSuggestion,
			}

			warning.Details = append(warning.Details, details)
		}

		warnings = append(warnings, warning)
	}

	return warnings, nil
}

func (e *Evaluator) parseEvaluationResultsToSummary(results *EvaluationResults, loginURL string) printer.Summary {

	totalRulesEvaluated := results.Summary.RulesCount * results.Summary.FilesCount
	plainRows := []printer.SummaryItem{
		{LeftCol: "Enabled rules in policy “default”", RightCol: fmt.Sprint(results.Summary.RulesCount), RowIndex: 0},
		{LeftCol: "Configs tested against policy", RightCol: fmt.Sprint(results.Summary.FilesCount), RowIndex: 1},
		{LeftCol: "Total rules evaluated", RightCol: fmt.Sprint(totalRulesEvaluated), RowIndex: 2},
		{LeftCol: "See all rules in policy", RightCol: loginURL, RowIndex: 5},
	}

	successRow := printer.SummaryItem{LeftCol: "Total rules passed", RightCol: fmt.Sprint(totalRulesEvaluated - results.Summary.TotalFailedRules), RowIndex: 4}
	errorRow := printer.SummaryItem{LeftCol: "Total rules failed", RightCol: fmt.Sprint(results.Summary.TotalFailedRules), RowIndex: 3}

	summary := &printer.Summary{
		ErrorRow:   errorRow,
		SuccessRow: successRow,
		PlainRows:  plainRows,
	}
	return *summary
}

type Rule struct {
	ID             int
	Name           string
	FailSuggestion string
	Count          int
}

func (rp *Rule) IncrementCount() {
	rp.Count++
}
