package vo

type Tag string

func DiffTags(a, b []Tag) (onlyA, onlyB []Tag) {
	aSet := make(map[Tag]struct{})
	bSet := make(map[Tag]struct{})

	// 将切片元素存入 map
	for _, item := range a {
		aSet[item] = struct{}{}
	}

	for _, item := range b {
		bSet[item] = struct{}{}
	}

	for e := range aSet {
		if _, ok := bSet[e]; !ok {
			onlyA = append(onlyA, e)
		}
	}

	for e := range bSet {
		if _, ok := aSet[e]; !ok {
			onlyB = append(onlyB, e)
		}
	}

	return onlyA, onlyB
}
