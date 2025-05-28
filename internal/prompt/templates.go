package prompt

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func LoadTemplates(tplPath string) ([]string, error) {
	if tplPath == "" {
		tplPath = filepath.Join(".", ".ai-agent-cli", "templates")
	}

	tplPath = strings.TrimPrefix(tplPath, "./")
	if !strings.HasPrefix(tplPath, ".ai-agent-cli") {
		tplPath = filepath.Join(".", ".ai-agent-cli", tplPath)
	}

	slog.Debug("loading templates...", slog.String("path", tplPath))

	info, err := os.Stat(tplPath)
	if err != nil || !info.IsDir() {
		return nil, nil // no templates, no return
	}

	var files []string

	if err := filepath.WalkDir(tplPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".tpl") {
			files = append(files, path)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error to read templates: %w", err)
	}

	sort.Strings(files)

	var contents []string

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("unable to read %s: %w", f, err)
		}

		header := fmt.Sprintf("### %s\n", filepath.Base(f))
		contents = append(contents, header+string(b))

		slog.Debug("template loaded", slog.String("template", filepath.Base(f)))
	}

	slog.Debug("templates loaded successfully")
	return contents, nil
}
