package convert

func FromDomain[A, M any, P fromPtr[A, M]](a A) (b M) {
	P(&b).FromDomain(a)
	return b
}

// fromPtr 约束类型必须是指向 B 的指针，且该指针类型实现了 fromDomain(A) 方法
type fromPtr[A any, M any] interface {
	*M

	FromDomain(A)
}

func SliceFromDomain[A, M any, Mp fromPtr[A, M]](as []A) (bs []M) {
	if len(as) == 0 {
		return make([]M, 0)
	}
	bs = make([]M, len(as))
	for i := range as {
		Mp(&bs[i]).FromDomain(as[i])
	}
	return bs
}

// intoPtr constraints *M can convert to A
// 'A' stands for Aggregate in domain, and 'M' stands for Model mapping to repo.
type intoPtr[M any, A any] interface {
	*M

	IntoDomain() (A, error)
}

func SliceIntoDomain[M, A any, P intoPtr[M, A]](ms []M) (as []A, err error) {
	if len(ms) == 0 {
		return nil, nil
	}

	as = make([]A, len(ms))
	for i := range as {
		if as[i], err = P(&ms[i]).IntoDomain(); err != nil {
			return nil, err
		}
	}
	return as, nil
}
