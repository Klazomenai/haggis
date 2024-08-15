package main

import (
  "haggis/cmd"
  "os"
)

// Entry point of the application
func main() {
  // If no arguments are provided, show usage
  if len(os.Args) < 2 {
    cmd.ShowUsage()
    os.Exit(1)
  }

  // Handle subcommands and flags
  cmd.Execute(os.Args)
}
