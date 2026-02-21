package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

var versionRegex = regexp.MustCompile(`const\s+AppVersion\s*=\s*"(\d+)\.(\d+)\.(\d+)"`)
var semverRegex = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

const versionFile = "internal/config/version.go"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func main() {
	args := os.Args[1:]
	bumpType := ""
	if len(args) > 0 {
		bumpType = args[0]
	}

	if bumpType != "" && bumpType != "patch" && bumpType != "minor" && bumpType != "major" && bumpType != "set" && bumpType != "current" {
		exitErr(errors.New("usage: go run ./scripts/release.go [patch|minor|major|set x.y.z|current]"))
	}

	major, minor, patch := getCurrentVersion()
	current := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	if bumpType == "current" {
		fmt.Println(current)
		return
	}

	setVersion := ""
	if bumpType == "set" {
		if len(args) < 2 {
			exitErr(errors.New("set requires a version argument, example: go run ./scripts/release.go set 1.2.0"))
		}
		setVersion = strings.TrimSpace(args[1])
		if !semverRegex.MatchString(setVersion) {
			exitErr(errors.New(`set format must be "x.y.z"`))
		}
	}

	checkGitStatusClean()

	newVersion := setVersion
	if newVersion == "" {
		newVersion = chooseVersion(bumpType, major, minor, patch)
	}
	if newVersion == "" {
		return
	}

	fmt.Printf("\n%s🚀 Bump version: v%s -> v%s%s\n\n", colorCyan, current, newVersion, colorReset)
	updateVersion(newVersion)
	commitAndPush(newVersion)

	fmt.Printf("\n%s✅ Done.%s\n", colorGreen, colorReset)
	fmt.Printf("%s📌 Pushed to origin/main.%s\n", colorGreen, colorReset)
}

func checkGitStatusClean() {
	fmt.Printf("\n%s🔍 [CHECK] git status, checking for uncommitted changes...%s\n", colorCyan, colorReset)

	statusUNO, err := runCmd("git", "status", "-uno")
	if err != nil {
		exitErr(fmt.Errorf("error git status -uno: %w", err))
	}
	if strings.Contains(statusUNO, "Your branch is behind") {
		exitErr(errors.New("release rejected: branch is behind remote. run git pull first"))
	}

	porcelain, err := runCmd("git", "status", "--porcelain")
	if err != nil {
		exitErr(fmt.Errorf("error git status --porcelain: %w", err))
	}

	porcelain = strings.TrimSpace(porcelain)
	if porcelain != "" {
		fmt.Printf("%s⚠️  Uncommitted/untracked files detected:%s\n", colorYellow, colorReset)
		for _, line := range strings.Split(porcelain, "\n") {
			clean := formatDirtyLine(line)
			if clean == "" {
				continue
			}
			fmt.Println(clean)
		}
		exitErr(errors.New("release aborted. working tree is dirty, please commit or stash your changes first"))
	}

	fmt.Printf("%s✅ Working tree is clean.%s\n", colorGreen, colorReset)
}

func getCurrentVersion() (int, int, int) {
	path := filepath.Clean(versionFile)
	content, err := os.ReadFile(path)
	if err != nil {
		exitErr(fmt.Errorf("error reading version file %s: %w", path, err))
	}

	m := versionRegex.FindSubmatch(content)
	if len(m) != 4 {
		exitErr(errors.New(`version format not found (const AppVersion = "x.y.z")`))
	}

	return mustAtoi(string(m[1])), mustAtoi(string(m[2])), mustAtoi(string(m[3]))
}

func updateVersion(newVersion string) {
	path := filepath.Clean(versionFile)
	content, err := os.ReadFile(path)
	if err != nil {
		exitErr(fmt.Errorf("error reading version file %s: %w", path, err))
	}

	replaced := regexp.MustCompile(`const\s+AppVersion\s*=\s*"[\d.]+"`).ReplaceAll(
		content,
		[]byte(fmt.Sprintf("const AppVersion = \"%s\"", newVersion)),
	)

	if string(replaced) == string(content) {
		exitErr(errors.New("release rejected: version unchanged"))
	}

	if err := os.WriteFile(path, replaced, 0o644); err != nil {
		exitErr(fmt.Errorf("error writing version file %s: %w", path, err))
	}

	fmt.Printf("%s✅ [OK] %s -> %s%s\n", colorGreen, path, newVersion, colorReset)
}

func commitAndPush(newVersion string) {
	commitMessage := fmt.Sprintf("chore(release): v%s 🚀", newVersion)

	if _, err := runCmd("git", "add", "."); err != nil {
		exitErr(fmt.Errorf("error git add: %w", err))
	}

	if _, err := runCmd("git", "commit", "-m", commitMessage); err != nil {
		exitErr(fmt.Errorf("error git commit: %w", err))
	}

	pushOut, err := runCmd("git", "push", "origin", "main")
	if err != nil {
		fmt.Printf("%s⚠️  Push failed. Rolling back release commit...%s\n", colorYellow, colorReset)
		rollbackOut, rollbackErr := runCmd("git", "reset", "--hard", "HEAD~1")
		if rollbackErr != nil {
			exitErr(fmt.Errorf("error git push origin main: %w\n%s\n\nRollback failed: %v\nRollback output: %s",
				err, strings.TrimSpace(pushOut), rollbackErr, strings.TrimSpace(rollbackOut)))
		}
		exitErr(fmt.Errorf("error git push origin main: %w\n%s\n\nRollback success: release commit reverted",
			err, strings.TrimSpace(pushOut)))
	}

	fmt.Printf("%s✅ [OK] commit created: %s%s\n", colorGreen, commitMessage, colorReset)
	fmt.Printf("%s✅ [OK] pushed: origin/main%s\n", colorGreen, colorReset)
}

func chooseVersion(bumpType string, major, minor, patch int) string {
	current := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	vPatch := fmt.Sprintf("%d.%d.%d", major, minor, patch+1)
	vMinor := fmt.Sprintf("%d.%d.%d", major, minor+1, 0)
	vMajor := fmt.Sprintf("%d.%d.%d", major+1, 0, 0)

	fmt.Printf("\n%s📦 Current version: v%s%s\n", colorCyan, current, colorReset)

	switch bumpType {
	case "patch":
		return vPatch
	case "minor":
		return vMinor
	case "major":
		return vMajor
	}

	items := []string{
		fmt.Sprintf("patch -> %s", vPatch),
		fmt.Sprintf("minor -> %s", vMinor),
		fmt.Sprintf("major -> %s", vMajor),
		"custom",
		"cancel",
	}

	prompt := promptui.Select{
		Label: "Choose bump type (use up/down arrows)",
		Items: items,
		Size:  5,
	}

	_, selected, err := prompt.Run()
	if err != nil {
		exitErr(fmt.Errorf("failed to choose version: %w", err))
	}

	switch selected {
	case items[0]:
		return vPatch
	case items[1]:
		return vMinor
	case items[2]:
		return vMajor
	case "custom":
		customPrompt := promptui.Prompt{
			Label: "Enter custom version (x.y.z)",
			Validate: func(input string) error {
				if !semverRegex.MatchString(strings.TrimSpace(input)) {
					return errors.New(`format must be "x.y.z"`)
				}
				return nil
			},
		}
		customVersion, err := customPrompt.Run()
		if err != nil {
			exitErr(fmt.Errorf("failed to read custom version: %w", err))
		}
		return strings.TrimSpace(customVersion)
	default:
		fmt.Printf("%s❌ Canceled.%s\n", colorYellow, colorReset)
		return ""
	}
}

func cleanPorcelainLine(line string) string {
	line = strings.TrimRight(line, "\n")
	if strings.TrimSpace(line) == "" {
		return ""
	}

	// Standard porcelain line format is usually: XY<space>PATH
	// Example: " M README.md", "M  file.go", "?? new_file.go"
	if len(line) >= 3 && line[2] == ' ' {
		return strings.TrimSpace(line[3:])
	}

	// Fallback for lines like: "M README.md" (single status + space + path)
	if len(line) >= 2 {
		return strings.TrimSpace(line[2:])
	}

	return line
}

func formatDirtyLine(line string) string {
	line = strings.TrimRight(line, "\n")
	if strings.TrimSpace(line) == "" {
		return ""
	}

	path := cleanPorcelainLine(line)
	if path == "" {
		return ""
	}

	return fmt.Sprintf("M       %s", path)
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func mustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		exitErr(fmt.Errorf("invalid version segment %q: %w", s, err))
	}
	return n
}

func exitErr(err error) {
	fmt.Fprintf(os.Stderr, "%s❌ [ERROR]: %v%s\n", colorRed, err, colorReset)
	os.Exit(1)
}
