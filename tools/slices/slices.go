package slices

func Equal[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func Clone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	result := make(S, len(s))
	copy(result, s)
	return result
}

func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	// Don't start copying elements until we find one to delete.
	for i, v := range s {
		if del(v) {
			j := i
			for i++; i < len(s); i++ {
				v = s[i]
				if !del(v) {
					s[j] = v
					j++
				}
			}
			return s[:j]
		}
	}
	return s
}
