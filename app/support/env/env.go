package env

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func IsAir() bool {
	for _, arg := range os.Args {
		if strings.Contains(filepath.ToSlash(arg), "/storage/temp") {
			return true
		}
	}

	return false
}

func IsDirectlyRun() bool {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Contains(filepath.Base(executable), os.TempDir()) ||
		(strings.Contains(filepath.ToSlash(executable), "/var/folders") &&
			strings.Contains(filepath.ToSlash(executable), "/T/go-build")) // macOS
}

func IsTesting() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "-test.") {
			return true
		}
	}

	return false
}

func CurrentAbsolutePath() string {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	res, err := filepath.EvalSymlinks(filepath.Dir(executable))
	if err != nil {
		log.Fatal(err)
	}

	if IsTesting() || IsAir() || IsDirectlyRun() {
		res, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	return res
}
