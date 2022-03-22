package common

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

var (
	DefaultDirMod   os.FileMode = 0755
	DefaultFileMode os.FileMode = 0600
)

// EnsurePath ensures a directory exist from the given path.
func EnsurePath(path string, mod os.FileMode) {
	dir := filepath.Dir(path)
	EnsureFullPath(dir, mod)
}

// EnsureFullPath ensures a directory exist from the given path.
func EnsureFullPath(path string, mod os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, mod); err != nil {
			log.Fatal().Msgf("Unable to create dir %q %v", path, err)
		}
	}
}

func Map[T any, V any](source []T, target []V, mapper func(T) V) []V {
	for _, t := range source {
		target = append(target, mapper(t))
	}
	return target
}
