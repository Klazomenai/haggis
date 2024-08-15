package cmd

import (
  "flag"
  "fmt"
  "haggis/internal"
  "os"
)

var (
  orgFlag     string
  branchFlag  string
  verboseFlag bool
)

// Execute handles the parsing of the command-line arguments and flags.
func Execute(args []string) {
  codeownersCmd := flag.NewFlagSet("codeowners", flag.ExitOnError)

  // Define flags for the codeowners subcommand
  codeownersCmd.StringVar(&orgFlag, "org", "", "GitHub organization to fetch CODEOWNERS from")
  codeownersCmd.StringVar(&branchFlag, "branch", "", "Specific branch to fetch the CODEOWNERS file from (optional)")
  codeownersCmd.BoolVar(&verboseFlag, "v", false, "Enable verbose output")

  // Parse the command-line arguments
  if len(args) < 2 {
    ShowUsage()
    os.Exit(1)
  }

  switch args[1] {
  case "codeowners":
    err := codeownersCmd.Parse(args[2:])
    if err != nil {
      fmt.Println("Error parsing codeowners flags")
      os.Exit(1)
    }

    // Check if the --org flag is provided
    if orgFlag == "" {
      fmt.Println("You must specify a GitHub organization with the --org flag.")
      ShowUsage()
      os.Exit(1)
    }

    // Fetch CODEOWNERS files for all repositories in the organization
    if err := internal.SearchCodeOwnersForOrg(orgFlag, branchFlag, verboseFlag); err != nil {
      fmt.Printf("Error: %v\n", err)
      os.Exit(1)
    }
  default:
    ShowUsage()
    os.Exit(1)
  }
}

// ShowUsage prints usage information for the application.
func ShowUsage() {
  fmt.Println("Usage: haggis <subcommand> [flags]")
  fmt.Println("Available subcommands:")
  fmt.Println("  codeowners - Process the CODEOWNERS file and handle GitHub org")
  fmt.Println("Flags:")
  fmt.Println("  --org <org>   : Specify the GitHub organization (required)")
  fmt.Println("  --branch <br> : Specify the branch to use (optional, default is the repo's default branch)")
  fmt.Println("  -v            : Enable verbose output")
}
