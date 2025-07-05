package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "specform",
		Short: "A CLI tool for compiling and validating prompt specification files.",
		Long:  `specform is a CLI tool for compiling and validating prompt specification files.`,
	}

	rootCmd.AddCommand(NewCompileCommand())
	rootCmd.AddCommand(NewRenderCommand())
	rootCmd.AddCommand(NewTestCommand())
	rootCmd.AddCommand(NewSnapshotCommand())
	rootCmd.AddCommand(NewServeCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
