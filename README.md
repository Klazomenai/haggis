### HAGGIS: Hitherto Another Github Graphql Inspection Search

<p align="center">
  <img src="https://github.com/klazomenai/haggis/blob/main/images/haggis-3.png?raw=true" alt="Haggis"/>
</p>

`HAGGIS` is a command-line tool built with Go, designed to search GitHub repositories using GraphQL. It simplifies the task of searching for specific files like `CODEOWNERS` across GitHub organizations and repositories, helping platform engineers streamline their searches.

### Installation

- Clone the repository:
```bash
git clone https://github.com/your-repo/haggis.git
cd haggis
```

- Build the binary:
```bash
go build -o haggis
```

- Run via Docker (optional):
```bash
docker build -t haggis .
```

### Usage

#### Command Structure

```bash
haggis [command] [flags]
```

#### Commands

- **`codeowners`**: Searches for the `.github/CODEOWNERS` file in the specified GitHub organization and repository.

#### Flags

- `--org, -o` (required): GitHub organization to search.
- `--repo, -r` (optional): Specific GitHub repository to search. If omitted, all repositories in the organization will be searched.
- `--branch, -b` (optional): Specific branch to search. If not specified, the default branch is used.

### Examples

- Search all repositories in an organization for the `CODEOWNERS` file:
```bash
haggis codeowners --org your-github-org
```

- Search a specific repository for the `CODEOWNERS` file on a specific branch:
```bash
haggis codeowners --org your-github-org --repo your-repo --branch main
```

- Run via Docker:
```bash
docker run --rm -it -e GITHUB_TOKEN=your_github_token --entrypoint haggis haggis codeowners --org your-github-org
```

### Notes

- Ensure the `GITHUB_TOKEN` environment variable is set with a valid GitHub Personal Access Token (PAT) for authenticated requests.
- HAGGIS supports auto-completion scripts for enhanced CLI experience using the `completion` sub-command.
