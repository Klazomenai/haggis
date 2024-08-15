package internal

import (
  "bufio"
  "fmt"
  "os"
  "strings"
)

// SearchCodeOwners reads and processes the CODEOWNERS file.
func SearchCodeOwners(file string, verbose bool) error {
  if verbose {
    fmt.Println(Colorize(fmt.Sprintf("Reading CODEOWNERS file: %s", file), Yellow))
  }

  // Open the CODEOWNERS file
  codeownersFile := file
  f, err := os.Open(codeownersFile)
  if err != nil {
    return fmt.Errorf("could not open CODEOWNERS file: %v", err)
  }
  defer f.Close()

  // Print the file header
  fmt.Println(Colorize("CODEOWNERS File:", Green))

  // Read the file line by line
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    line := scanner.Text()

    // Skip empty lines and comments
    if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
      continue
    }

    // Print valid entries
    fmt.Println(Colorize(line, Blue))
  }

  // Handle any errors during file reading
  if err := scanner.Err(); err != nil {
    return fmt.Errorf("error reading CODEOWNERS file: %v", err)
  }

  return nil
}
