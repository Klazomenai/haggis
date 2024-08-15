package internal

import (
  "fmt"
  "os"
  "strings"
)

// CalculateDefaultBranch determines the default branch by reading the .git/HEAD file.
func CalculateDefaultBranch(verbose bool) (string, error) {
  headFile := ".git/HEAD"
  if verbose {
    fmt.Println(Colorize(fmt.Sprintf("Reading default branch from: %s", headFile), Yellow))
  }

  // Open the HEAD file to determine the branch
  file, err := os.Open(headFile)
  if err != nil {
    return "", fmt.Errorf("could not open HEAD file: %v", err)
  }
  defer file.Close()

  // Read the branch reference
  buffer := make([]byte, 1024)
  n, err := file.Read(buffer)
  if err != nil {
    return "", fmt.Errorf("could not read HEAD file: %v", err)
  }

  content := string(buffer[:n])
  if strings.HasPrefix(content, "ref: ") {
    // Extract the branch name
    parts := strings.Split(content, "/")
    branch := strings.TrimSpace(parts[len(parts)-1])
    if verbose {
      fmt.Println(Colorize(fmt.Sprintf("Default branch detected: %s", branch), Blue))
    }
    return branch, nil
  }

  // Handle detached head (commit hash)
  return strings.TrimSpace(content), nil
}
