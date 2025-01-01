package seq

import "iter"

func Map[U, V any, F ~func(U) V](seq Seq[U], mapper F) Seq[V] {
	return func() (V, bool) {
		value, valid := seq()
		return mapper(value), valid
	}
}

func MapSeq[U, V any, F ~func(U) V](seq iter.Seq[U], mapper F) iter.Seq[V] {
	return func(yield func(V) bool) {
		for value := range seq {
			if !yield(mapper(value)) {
				break
			}
		}
	}
}

func Map2[K, V, MK, MV any, F ~func(K, V) (MK, MV)](seq Seq2[K, V], mapper F) Seq2[MK, MV] {
	return func() (MK, MV, bool) {
		key, value, valid := seq()
		mk, mv := mapper(key, value)
		return mk, mv, valid
	}
}

func MapSeq2[K, V, MK, MV any, F ~func(K, V) (MK, MV)](seq iter.Seq2[K, V], mapper F) iter.Seq2[MK, MV] {
	return func(yield func(MK, MV) bool) {
		for key, value := range seq {
			mk, mv := mapper(key, value)
			if !yield(mk, mv) {
				break
			}
		}
	}
}

func FlatMap[U, V any, F ~func(U) []V](seq Seq[U], mapper F) Seq[V] {
	remaining := make(chan V, 100)
	return func() (V, bool) {
		for len(remaining) == 0 {
			value, valid := seq()
			if valid {
				mapped := mapper(value)
				for _, v := range mapped {
					remaining <- v
				}
			} else {
				break
			}
		}
		for len(remaining) > 0 {
			return <-remaining, true
		}
		close(remaining)
		return *new(V), false
	}
}

func FlatMapSeq[U, V any, F ~func(U) []V](seq iter.Seq[U], mapper F) iter.Seq[V] {
	remaining := make(chan V, 100)
	return func(yield func(V) bool) {
		for value := range seq {
			mapped := mapper(value)
			for _, v := range mapped {
				remaining <- v
			}
			if len(remaining) > 0 {
				if !yield(<-remaining) {
					break
				}
			}
		}
		for len(remaining) > 0 {
			if !yield(<-remaining) {
				break
			}
		}
		close(remaining)
	}
}

func FlatMap2[K, V, MK, MV any, F ~func(K, V) []Pair[MK, MV]](seq Seq2[K, V], mapper F) Seq2[MK, MV] {
	remaining := make(chan Pair[MK, MV], 100)
	return func() (MK, MV, bool) {
		for len(remaining) == 0 {
			key, value, valid := seq()
			if valid {
				mapped := mapper(key, value)
				for _, v := range mapped {
					remaining <- v
				}
			} else {
				break
			}
		}
		for len(remaining) > 0 {
			kv := <-remaining
			return kv.A, kv.B, true
		}
		close(remaining)
		return *new(MK), *new(MV), false
	}
}

func FlatMapSeq2[K, V, MK, MV any, F ~func(K, V) []Pair[MK, MV]](seq iter.Seq2[K, V], mapper F) iter.Seq2[MK, MV] {
	remaining := make(chan Pair[MK, MV], 100)
	return func(yield func(MK, MV) bool) {
		for key, value := range seq {
			mapped := mapper(key, value)
			for _, kv := range mapped {
				remaining <- kv
			}
			if len(remaining) > 0 {
				kv := <-remaining
				if !yield(kv.A, kv.B) {
					break
				}
			}
		}
		for len(remaining) > 0 {
			kv := <-remaining
			if !yield(kv.A, kv.B) {
				break
			}
		}
		close(remaining)
	}
}
