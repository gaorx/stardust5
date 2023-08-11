package sdsemver

func Equal(a, b V) bool {
	return a == b
}

func Compare(a, b V) int {
	if a.Major < b.Major {
		return -1
	} else if a.Major > b.Major {
		return 1
	}
	if a.Minor < b.Minor {
		return -1
	} else if a.Minor > b.Minor {
		return 1
	}
	if a.Patch < b.Patch {
		return -1
	} else if a.Patch > b.Patch {
		return 1
	}
	return 0
}
