package main

// system represents a supported system/arch
type system struct {
	name  string
	archs []string
}

// HasArch loops the archs list and returns if a match is found.
func (s *system) HasArch(arch string) bool {
	for _, a := range s.archs {
		if arch == a {
			return true
		}
	}
	return false
}

// systemList is a simple wrapper for a slice of system pointers.
type systemList []*system

// Find will return the system with the matching os name.
func (list systemList) Find(os string) *system {
	for _, s := range list {
		if s.name == os {
			return s
		}
	}
	return nil
}
