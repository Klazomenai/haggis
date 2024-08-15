package main

import (
  "context"
  "fmt"
  "os"
  "strings"

  "github.com/fatih/color"
  "github.com/shurcooL/githubv4"
  "github.com/spf13/cobra"
  "golang.org/x/oauth2"
)

type repositoryQuery struct {
  Repository struct {
    Object struct {
      Blob struct {
        Text string
      } `graphql:"... on Blob"`
    } `graphql:"object(expression: $expression)"`
  } `graphql:"repository(owner: $owner, name: $name)"`
}

func main() {
  // Create the root command for the CLI tool
  var rootCmd = &cobra.Command{
    Use:   "haggis",
    Short: "Haggis: Search GitHub repositories for specific files and configurations",
  }

  // Variables to hold flag values
  var org, repo, branch string

  // Create the codeowners subcommand
  var codeownersCmd = &cobra.Command{
    Use:   "codeowners",
    Short: "Search for the .github/CODEOWNERS file in the specified organization and repository",
    Run: func(cmd *cobra.Command, args []string) {
      if org == "" {
        color.Red("Error: Organization name is required.")
        return
      }

      // Search for the CODEOWNERS file
      if err := searchCodeOwners(org, repo, branch); err != nil {
        color.Red("Error: %v\n", err)
      }
    },
  }

  // Add flags to the codeowners command
  codeownersCmd.Flags().StringVarP(&org, "org", "o", "", "GitHub organization (required)")
  codeownersCmd.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository (optional)")
  codeownersCmd.Flags().StringVarP(&branch, "branch", "b", "", "GitHub branch (optional)")

  // Add the codeowners command as a subcommand of the root command
  rootCmd.AddCommand(codeownersCmd)

  // Execute the root command
  if err := rootCmd.Execute(); err != nil {
    color.Red("Error: %v\n", err)
    os.Exit(1)
  }
}

func searchCodeOwners(org, repo, branch string) error {
    src := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
    )
    httpClient := oauth2.NewClient(context.Background(), src)

    client := githubv4.NewClient(httpClient)

    // If no repo specified, search all repos in the organization
    if repo == "" {
        repos, err := listOrgRepos(client, org)
        if err != nil {
            return err
        }

        for _, r := range repos {
            // Determine the branch (use default if not provided)
            if branch == "" {
                defaultBranch, err := getDefaultBranch(client, org, r)
                if err != nil {
                    color.Red("Error getting default branch for repository %s: %v\n", r, err)
                    continue
                }
                branch = defaultBranch
            }

            // Now print the message with the correct branch
            color.Cyan("Searching CODEOWNERS in organization '%s', repository '%s', branch '%s'\n", org, r, branch)
            if err := searchRepoForCodeowners(client, org, r, branch); err != nil {
                color.Red("Error searching in repository %s: %v\n", r, err)
            }

            // Reset the branch to empty string for next repository (if iterating multiple repos)
            branch = ""
        }
        return nil
    }

    // If a specific repo is provided
    if branch == "" {
        // Determine the branch (use default if not provided)
        defaultBranch, err := getDefaultBranch(client, org, repo)
        if err != nil {
            return err
        }
        branch = defaultBranch
    }

    // Now print the message with the correct branch
    color.Cyan("Searching CODEOWNERS in organization '%s', repository '%s', branch '%s'\n", org, repo, branch)
    return searchRepoForCodeowners(client, org, repo, branch)
}

func listOrgRepos(client *githubv4.Client, org string) ([]string, error) {
  // GraphQL query to list repositories in the organization
  var query struct {
    Organization struct {
      Repositories struct {
        Nodes []struct {
          Name string
        }
      } `graphql:"repositories(first: 100)"`
    } `graphql:"organization(login: $login)"`
  }

  variables := map[string]interface{}{
    "login": githubv4.String(org),
  }

  err := client.Query(context.Background(), &query, variables)
  if err != nil {
    return nil, err
  }

  var repos []string
  for _, repo := range query.Organization.Repositories.Nodes {
    repos = append(repos, repo.Name)
  }
  return repos, nil
}

func searchRepoForCodeowners(client *githubv4.Client, owner, repo, branch string) error {
  // Determine the branch (use default if not provided)
  if branch == "" {
    defaultBranch, err := getDefaultBranch(client, owner, repo)
    if err != nil {
      return err
    }
    branch = defaultBranch
  }

  // Query the .github/CODEOWNERS file in the repository
  var query repositoryQuery
  expression := fmt.Sprintf("%s:.github/CODEOWNERS", branch)
  variables := map[string]interface{}{
    "owner":      githubv4.String(owner),
    "name":       githubv4.String(repo),
    "expression": githubv4.String(expression),
  }

  err := client.Query(context.Background(), &query, variables)
  if err != nil {
    return err
  }

  // If CODEOWNERS file exists, print its contents excluding comments
  if query.Repository.Object.Blob.Text != "" {
    content := query.Repository.Object.Blob.Text
    printCodeownersContent(content)
  } else {
    color.Yellow("No CODEOWNERS file found in %s/%s on branch %s\n", owner, repo, branch)
  }

  return nil
}

func getDefaultBranch(client *githubv4.Client, owner, repo string) (string, error) {
  var query struct {
    Repository struct {
      DefaultBranchRef struct {
        Name string
      }
    } `graphql:"repository(owner: $owner, name: $name)"`
  }

  variables := map[string]interface{}{
    "owner": githubv4.String(owner),
    "name":  githubv4.String(repo),
  }

  err := client.Query(context.Background(), &query, variables)
  if err != nil {
    return "", err
  }

  return query.Repository.DefaultBranchRef.Name, nil
}

func printCodeownersContent(content string) {
  lines := strings.Split(content, "\n")
  for _, line := range lines {
    line = strings.TrimSpace(line)
    if len(line) == 0 || strings.HasPrefix(line, "#") {
      continue
    }
    color.Green(line)
  }
}

