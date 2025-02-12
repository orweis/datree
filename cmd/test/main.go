package test

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/datreeio/datree/bl"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/propertiesExtractor"
	"github.com/spf13/cobra"
)

type LocalConfigManager interface {
	GetConfiguration() (localConfig.LocalConfiguration, error)
}

type Evaluator interface {
	PrintResults(results *bl.EvaluationResults, cliId string, output string) error
	PrintFileParsingErrors(errors []propertiesExtractor.FileError)
	Evaluate(pattern string, cliId string, evaluationConc int) (*bl.EvaluationResults, []propertiesExtractor.FileError, error)
}

type TestCommandContext struct {
	LocalConfig LocalConfigManager
	Evaluator   Evaluator
}

type TestCommandFlags struct {
	Output string
}

func CreateTestCommand(ctx *TestCommandContext) *cobra.Command {
	testCommand := &cobra.Command{
		Use:   "test",
		Short: "Execute static analysis for pattern",
		Long:  `Execute static analysis for pattern. Input should be glob`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFlag, err := cmd.Flags().GetString("output")
			if err != nil {
				fmt.Println(err)
				return err
			}

			testCommandFlags := TestCommandFlags{Output: outputFlag}
			return test(ctx, args[0], testCommandFlags)
		},
		Args:          cobra.ExactValidArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	testCommand.Flags().StringP("output", "o", "", "Define output format")
	return testCommand
}

func test(ctx *TestCommandContext, pattern string, flags TestCommandFlags) error {
	absolutePath, err := filepath.Abs(pattern)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	s.Suffix = " Loading..."
	s.Color("cyan")
	s.Start()

	config, err := ctx.LocalConfig.GetConfiguration()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	evaluationResponse, fileParsingErrors, err := ctx.Evaluator.Evaluate(absolutePath, config.CliId, 50)
	s.Stop()

	if err != nil {
		if len(fileParsingErrors) > 0 {
			ctx.Evaluator.PrintFileParsingErrors(fileParsingErrors)
		}
		fmt.Println(err.Error())
		return err
	}

	if evaluationResponse == nil {
		err := fmt.Errorf("no response received")
		return err

	}

	return ctx.Evaluator.PrintResults(evaluationResponse, config.CliId, flags.Output)
}
