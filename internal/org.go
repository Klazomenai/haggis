package internal

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "strings"
)

// Repo represents a GitHub repository (simplified)
type Repo struct {
  Name          string `json:"name"`
  DefaultBranch string `json:"default_branch"`
}

// FetchRepositories fetches all repositories for a GitHub organization.
func FetchRepositories(org string) ([]Repo, error) {
  reposURL := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
  resp, err := http.Get(reposURL)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch repositories for organization: %v", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    return nil, fmt.Errorf("failed to fetch repositories: received status code %d", resp.StatusCode)
  }

  // Read and parse the repositories JSON
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read repositories response: %v", err)
  }

  var repos []Repo
  err = json.Unmarshal(body, &repos)
  if err != nil {
    return nil, fmt.Errorf("failed to parse repositories data: %v", err)
  }

  return repos, nil
}

// GetCodeownersFile fetches the CODEOWNERS file from the specified branch of a repository.
func GetCodeownersFile(org, repo, branch string, verbose bool) (string, error) {
  if verbose {
    fmt.Println(Colorize(fmt.Sprintf("Fetching CODEOWNERS from branch '%s' for repository '%s'", branch, repo), Yellow))
  }

  // Construct the URL to fetch the CODEOWNERS file from the specified branch
  codeownersURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/.github/CODEOWNERS", org, repo, branch)
  resp, err := http.Get(codeownersURL)
  if err != nil {
    return "", fmt.Errorf("failed to fetch CODEOWNERS for repository %s: %v", repo, err)
  }
  defer resp.Body.Close()

  // If no CODEOWNERS file found, return empty
  if resp.StatusCode != 200 {
    return "", nil
  }

  // Read the CODEOWNERS file content
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", fmt.Errorf("failed to read CODEOWNERS file for repository %s: %v", repo, err)
  }

  // Return the content as a string, stripping out comments and empty lines
  content := string(body)
  codeowners := filterCodeownersContent(content)
  return codeowners, nil
}

// filterCodeownersContent removes comments and empty lines from the CODEOWNERS content.
func filterCodeownersContent(content string) string {
  var filteredLines []string
  lines := strings.Split(content, "\n")
  for _, line := range lines {
    trimmed := strings.TrimSpace(line)
    if trimmed == "" || strings.HasPrefix(trimmed, "#") {
      continue
    }
    filteredLines = append(filteredLines, trimmed)
  }
  return strings.Join(filteredLines, "\n")
}

// SearchCodeOwnersForOrg fetches repositories from a GitHub organization and checks for CODEOWNERS.
func SearchCodeOwnersForOrg(org string, branch string, verbose bool) error {
  repos, err := FetchRepositories(org)
  if err != nil {
    return err
  }

  foundAny := false
  for _, repo := range repos {
    // Use the default branch unless the user has specified a branch
    branchToUse := branch
    if branch == "" {
      branchToUse = repo.DefaultBranch
    }

    // Fetch the CODEOWNERS file
    codeowners, err := GetCodeownersFile(org, repo.Name, branchToUse, verbose)
    if err != nil {
      fmt.Printf("Failed to fetch CODEOWNERS for repository %s: %v\n", repo.Name, err)
      continue
    }

    if codeowners != "" {
      fmt.Println(Colorize(fmt.Sprintf("CODEOWNERS file found in repository '%s' (branch: %s):", repo.Name, branchToUse), Green))
      fmt.Println(codeowners)
      foundAny = true
    } else {
      if verbose {
        fmt.Println(fmt.Sprintf("No CODEOWNERS file found in repository '%s' (branch: %s)", repo.Name, branchToUse))
      }
    }
  }

  if !foundAny {
    return fmt.Errorf("no CODEOWNERS files found for organization %s", org)
  }

  return nil
}
