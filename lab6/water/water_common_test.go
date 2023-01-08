package water

func isWaterMolecule(mol string) bool {
	hCount := 0
	oCount := 0
	for _, ch := range mol {
		switch ch {
		case 'H':
			hCount++
		case 'O':
			oCount++
		}
	}
	return hCount == 2 && oCount == 1
}

func isWaterSequence(s string) bool {
	if len(s) == 0 || len(s)%3 != 0 {
		return false
	}
	for i := 0; i < len(s); i += 3 {
		if !isWaterMolecule(s[i : i+3]) {
			return false
		}
	}
	return true
}
