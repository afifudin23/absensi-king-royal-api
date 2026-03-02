package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var releaseHeaderRe = regexp.MustCompile(`(?m)^## \[([^\]]+)\] - .*$`)
var appVersionRe = regexp.MustCompile(`(?m)^const\s+AppVersion\s*=\s*"([^"]+)"\s*$`)

func main() {
	var (
		tagName   = flag.String("tag", "", "git tag name override (default: v<AppVersion>)")
		latest    = flag.Bool("latest", false, "also create/update 'latest' tag")
		printOnly = flag.Bool("print", false, "print generated tag message only")
		noPush    = flag.Bool("no-push", false, "do not push tag(s) to origin")
	)
	flag.Parse()

	appVersion, err := readAppVersion("internal/config/version.go")
	if err != nil {
		exitErr(err)
	}

	body, err := extractReleaseBodyByVersion("CHANGELOG.md", appVersion)
	if err != nil {
		exitErr(err)
	}

	targetTag := strings.TrimSpace(*tagName)
	if targetTag == "" {
		targetTag = "v" + appVersion
	}

	message := fmt.Sprintf("Release v%s\n\n%s", appVersion, body)

	if *printOnly {
		fmt.Println(message)
		return
	}

	if err := runGitTag(targetTag, message); err != nil {
		exitErr(err)
	}
	fmt.Printf("✅ tag '%s' updated from latest changelog (v%s)\n", targetTag, appVersion)

	if *latest {
		if err := runGitTag("latest", message); err != nil {
			exitErr(err)
		}
		fmt.Printf("✅ tag 'latest' updated from latest changelog (v%s)\n", appVersion)
	}

	if !*noPush {
		if _, err := runCmd("git", "push", "origin", targetTag, "--force"); err != nil {
			exitErr(fmt.Errorf("failed to push tag '%s': %w", targetTag, err))
		}
		fmt.Printf("✅ tag '%s' pushed to origin\n", targetTag)

		if *latest {
			if _, err := runCmd("git", "push", "origin", "latest", "--force"); err != nil {
				exitErr(fmt.Errorf("failed to push tag 'latest': %w", err))
			}
			fmt.Println("✅ tag 'latest' pushed to origin")
		}
	}
}

func readAppVersion(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	m := appVersionRe.FindStringSubmatch(string(raw))
	if len(m) != 2 {
		return "", fmt.Errorf("AppVersion not found in %s", path)
	}
	v := strings.TrimSpace(m[1])
	if v == "" {
		return "", fmt.Errorf("AppVersion empty in %s", path)
	}
	return v, nil
}

func extractReleaseBodyByVersion(path string, version string) (body string, err error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}

	content := string(raw)
	indices := releaseHeaderRe.FindAllStringSubmatchIndex(content, -1)
	if len(indices) == 0 {
		return "", errors.New("no release section found in CHANGELOG.md")
	}

	start := -1
	end := -1
	for i, idx := range indices {
		v := strings.TrimSpace(content[idx[2]:idx[3]])
		if v != version {
			continue
		}
		start = idx[0]
		end = len(content)
		if i+1 < len(indices) {
			end = indices[i+1][0]
		}
		break
	}
	if start == -1 || end == -1 {
		return "", fmt.Errorf("release section for version %s not found in CHANGELOG.md", version)
	}

	section := strings.TrimSpace(content[start:end])
	lines := strings.Split(section, "\n")
	if len(lines) == 0 {
		return "", errors.New("release changelog section is empty")
	}

	// Remove top header (## [x.y.z] - yyyy-mm-dd).
	if strings.HasPrefix(strings.TrimSpace(lines[0]), "## ") {
		lines = lines[1:]
	}

	lines = removeMigrationSection(lines)
	lines = removeMigrationBullets(lines)

	body = strings.TrimSpace(strings.Join(lines, "\n"))
	if body == "" {
		return "", errors.New("release changelog body is empty after filtering")
	}

	return body, nil
}

func removeMigrationSection(lines []string) []string {
	out := make([]string, 0, len(lines))
	skip := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### Migration Required") {
			skip = true
			continue
		}

		if skip && strings.HasPrefix(trimmed, "### ") {
			skip = false
		}

		if !skip {
			out = append(out, line)
		}
	}

	// Remove trailing separators if present.
	for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "---" {
		out = out[:len(out)-1]
	}

	return out
}

func removeMigrationBullets(lines []string) []string {
	out := make([]string, 0, len(lines))
	skip := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lowerTrimmed := strings.ToLower(trimmed)
		isTopLevelBullet := strings.HasPrefix(line, "- ")

		// Skip bullet blocks like "- Attendance migration:" and nested bullets below it.
		if isTopLevelBullet && strings.Contains(lowerTrimmed, "migration") {
			skip = true
			continue
		}

		if skip {
			// End skipping when next top-level bullet starts and it's not an indented sub-bullet.
			if isTopLevelBullet {
				skip = false
				// Continue processing this new top-level bullet.
			} else {
				continue
			}
		}

		out = append(out, line)
	}

	return out
}

func runGitTag(tagName, message string) error {
	if strings.TrimSpace(tagName) == "" {
		return errors.New("tag name cannot be empty")
	}

	_, err := runCmd("git", "tag", "-fa", tagName, "-m", message)
	if err != nil {
		return fmt.Errorf("failed to create/update tag '%s': %w", tagName, err)
	}
	return nil
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%w (%s)", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func exitErr(err error) {
	fmt.Fprintf(os.Stderr, "❌ %v\n", err)
	os.Exit(1)
}
