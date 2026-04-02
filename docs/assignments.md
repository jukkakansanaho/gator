# Gator guidelines
This file contains assignments (step-by-step development guidelines) for developing Gator program.

# Assignments

## Assignment 1

### Step 1 – Create the config file

Create the file `~/.gatorconfig.json` in your home directory with the following content:

```json
{
  "db_url": "postgres://example"
}
```

Leave out `current_user_name` for now — the application will set it automatically.

### Step 2 – Initialize the Go module and main entry point

Initialize a new Go module (use any module path you like, e.g. `github.com/<username>/gator`):

```bash
go mod init github.com/<username>/gator
```

Create `main.go` at the project root with an empty `main` function:

```go
package main

func main() {}
```

### Step 3 – Create the internal config package

Create the directory structure for the internal config package:

```
internal/
└── config/
    └── config.go
```

The `internal` directory prevents other Go modules from importing these packages directly.

### Step 4 – Implement the config package

In `internal/config/config.go`, implement the following:

**Constant** — holds the config filename:

```go
const configFileName = ".gatorconfig.json"
```

**Struct** — represents the JSON file structure, with JSON struct tags:

```go
type Config struct {
    DbURL           string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}
```

**Unexported helpers:**

- `getConfigFilePath() (string, error)` — uses `os.UserHomeDir()` to build the full path to `~/.gatorconfig.json`
- `write(cfg Config) error` — marshals the `Config` struct to JSON and writes it to the config file path

**Exported functions/methods:**

- `Read() (Config, error)` — calls `getConfigFilePath`, opens the file, and decodes the JSON into a new `Config` struct
- `(c *Config) SetUser(name string) error` — sets `c.CurrentUserName = name` then calls `write(*c)` to persist the change

### Step 5 – Wire up main

Update `main.go` to exercise the config package:

1. Call `config.Read()` to load the config from disk.
2. Call `cfg.SetUser("your-name")` to set the current user and write the updated config back to disk.
3. Call `config.Read()` a second time and print the resulting struct to the terminal with `fmt.Printf` or `fmt.Println`.

After running `go run .`, you should see the config struct printed with both `db_url` and `current_user_name` populated, and `~/.gatorconfig.json` updated on disk.

### Step 6 - Write tests

Create `internal/config/config_test.go` in the same package (`package config`) so tests have access to unexported helpers.

#### Test isolation strategy

`Read`, `write`, and `getConfigFilePath` all resolve to `~/.gatorconfig.json` via `os.UserHomeDir`. To avoid touching the real config file, redirect `HOME` to a temporary directory for each test using `t.Setenv` and `t.TempDir`:

```go
func useHome(t *testing.T, dir string) {
    t.Helper()
    t.Setenv("HOME", dir)
}
```

`t.Setenv` automatically restores the original value when the test ends.

Add a helper that writes a fixture config file into a given directory:

```go
func writeFixture(t *testing.T, dir string, cfg Config) {
    t.Helper()
    data, _ := json.Marshal(cfg)
    os.WriteFile(filepath.Join(dir, configFileName), data, 0600)
}
```

#### Tests to write

- **`TestGetConfigFilePath`** — set `HOME` to a temp dir, call `getConfigFilePath()`, assert the returned path equals `filepath.Join(tmp, configFileName)`.

- **`TestRead_ValidFile`** — write a fixture config, call `Read()`, assert the returned `Config` matches the fixture.

- **`TestRead_MissingFile`** — use an empty temp dir (no config file), call `Read()`, assert an error is returned.

- **`TestRead_InvalidJSON`** — write a file containing `"not json"`, call `Read()`, assert an error is returned.

- **`TestSetUser`** — write a fixture config, call `Read()`, call `cfg.SetUser("name")`, call `Read()` again, assert `CurrentUserName` is updated and `DbURL` is unchanged.

- **`TestWrite`** — call `write(cfg)` directly, read the raw file back with `os.ReadFile`, unmarshal it, and assert the result matches the original `Config`.

#### Running the tests

```bash
go test ./internal/config/ -v
```

All tests should pass without modifying `~/.gatorconfig.json`.
