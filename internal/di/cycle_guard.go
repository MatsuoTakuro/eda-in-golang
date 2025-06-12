package di

import (
	"maps"
	"strings"
)

// cycleGuard keeps track of dependencies currently being constructed.
// It is used to detect and prevent cyclic dependencies at runtime.
type cycleGuard map[Key]int

// add returns a new tracked map with the given depInfo added.
func (cg cycleGuard) add(dep depDef) cycleGuard {
	dst := make(cycleGuard, len(cg))

	// Copy existing keys and their indices to the new list
	maps.Copy(dst, cg)
	dst[dep.key] = len(dst)

	return dst
}

func (cg cycleGuard) keys() []string {
	keys := make([]string, len(cg))

	for key, i := range cg {
		keys[i] = string(key)
	}

	return keys
}

func (cg cycleGuard) String() string {
	return strings.Join(cg.keys(), ",")
}
