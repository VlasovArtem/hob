package common

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
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

func MapKeys[T comparable, V comparable](source map[T]any, mapper func(T) V) map[V]any {
	if source == nil {
		return nil
	}
	if len(source) == 0 {
		return make(map[V]any)
	}
	result := make(map[V]any)
	for key, value := range source {
		result[mapper(key)] = value
	}
	return result
}

func MapValues[KEY comparable, T any, V any](source map[KEY]T, mapper func(T) V) map[KEY]V {
	if source == nil {
		return nil
	}
	result := make(map[KEY]V)
	if len(source) == 0 {
		return result
	}
	for key, value := range source {
		result[key] = mapper(value)
	}
	return result
}

func MapData[OKEY comparable, NKEY comparable, OVALUE any, NVALUE any](source map[OKEY]OVALUE, mapper func(OKEY, OVALUE) (NKEY, NVALUE)) map[NKEY]NVALUE {
	if source == nil {
		return nil
	}
	result := make(map[NKEY]NVALUE)
	if len(source) == 0 {
		return result
	}
	for key, value := range source {
		nkey, nvalue := mapper(key, value)
		result[nkey] = nvalue
	}
	return result
}

func MapSlice[T any, V any](source []T, mapper func(T) V) []V {
	target := make([]V, 0)
	for _, t := range source {
		target = append(target, mapper(t))
	}
	return target
}

func ForEach[T any](source []T, action func(T)) {
	for _, t := range source {
		action(t)
	}
}

func Join[T fmt.Stringer](t []T, sep string) string {
	return strings.Join(MapSlice(t, func(t T) string {
		return t.String()
	}), sep)
}
