package cmd

import (
  "flag"
  "fmt"
  "haggis/internal"
  "os"
)

// ExecuteCodeowners processes the codeowners subcommand and related flags.
func ExecuteCodeowners(args []string) {
  // Define the flags for this subcommand
  var fileFlag string
  var verboseFlag bool
  codeownersCmd := flag.NewFlagSet("codeowners", flag.ExitOnError)
  codeownersCmd.StringVar(&fileFlag, "f", ".github/CODEOWNERS", "Path to CODEOWNERS file")
  codeownersCmd.BoolVar(&verboseFlag, "v", false, "Enable verbose output")

  // Parse the flags for this subcommand
  err := codeownersCmd.Parse(args)
  if err != nil {
    fmt.Println("Error parsing codeowners flags")
    os.Exit(1)
  }

  // Execute the CODEOWNERS search logic, using the file and verbosity flags
  if err := internal.SearchCodeOwners(fileFlag, verboseFlag); err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
  }
}
