package seq

import "iter"

func Map[U, V any](seq Seq[U], mapper MapFunc[U, V]) Seq[V] {
	return func() (V, bool) {
		value, valid := seq()
		return mapper(value), valid
	}
}

func MapSeq[U, V any](seq iter.Seq[U], mapper MapFunc[U, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for value := range seq {
			if !yield(mapper(value)) {
				break
			}
		}
	}
}

func Map2[K, V, MK, MV any](seq Seq2[K, V], mapper Map2Func[K, V, MK, MV]) Seq2[MK, MV] {
	return func() (MK, MV, bool) {
		key, value, valid := seq()
		mk, mv := mapper(key, value)
		return mk, mv, valid
	}
}

func MapSeq2[K, V, MK, MV any](seq iter.Seq2[K, V], mapper Map2Func[K, V, MK, MV]) iter.Seq2[MK, MV] {
	return func(yield func(MK, MV) bool) {
		for key, value := range seq {
			mk, mv := mapper(key, value)
			if !yield(mk, mv) {
				break
			}
		}
	}
}

func FlatMap[U, V any](seq Seq[U], mapper FlatMapFunc[U, V]) Seq[V] {
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
		if len(remaining) > 0 {
			return <-remaining, true
		}
		close(remaining)
		return *new(V), false
	}
}

func FlatMapSeq[U, V any](seq iter.Seq[U], mapper FlatMapFunc[U, V]) iter.Seq[V] {
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
		close(remaining)
	}
}

func FlatMap2[K, V, MK, MV any](seq Seq2[K, V], mapper FlatMap2Func[K, V, MK, MV]) Seq2[MK, MV] {
	remaining := make(chan KeyValue[MK, MV], 100)
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
		if len(remaining) > 0 {
			kv := <-remaining
			return kv.Key, kv.Value, true
		}
		close(remaining)
		return *new(MK), *new(MV), false
	}
}

func FlatMapSeq2[K, V, MK, MV any](seq iter.Seq2[K, V], mapper FlatMap2Func[K, V, MK, MV]) iter.Seq2[MK, MV] {
	remaining := make(chan KeyValue[MK, MV], 100)
	return func(yield func(MK, MV) bool) {
		for key, value := range seq {
			mapped := mapper(key, value)
			for _, kv := range mapped {
				remaining <- kv
			}
			if len(remaining) > 0 {
				kv := <-remaining
				if !yield(kv.Key, kv.Value) {
					break
				}
			}
		}
		close(remaining)
	}
}
