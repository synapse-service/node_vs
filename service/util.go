package service

func Separate(a, b []string) ([]string, []string, []string) {
	am := make(map[string]struct{}, len(a))
	for _, v := range a {
		am[v] = struct{}{}
	}

	bm := make(map[string]struct{}, len(b))
	for _, v := range b {
		bm[v] = struct{}{}
	}

	inm := map[string]struct{}{}

	for k := range am {
		if _, ok := bm[k]; ok {
			inm[k] = struct{}{}
			delete(am, k)
			delete(bm, k)
		}
	}

	a1 := make([]string, 0, len(am))
	for k := range am {
		a1 = append(a1, k)
	}

	in := make([]string, 0, len(inm))
	for k := range inm {
		in = append(in, k)
	}

	b1 := make([]string, 0, len(bm))
	for k := range bm {
		b1 = append(b1, k)
	}

	return a1, in, b1
}
