package seq

func Map[U, V any](seq Seq[U], mapper MapFunc[U, V]) Seq[V] {
	return func() (V, bool) {
		value, valid := seq()
		return mapper(value), valid
	}
}

func Map2[K, V, MK, MV any](seq Seq2[K, V], mapper Map2Func[K, V, MK, MV]) Seq2[MK, MV] {
	return func() (MK, MV, bool) {
		for {
			key, value, valid := seq()
			mk, mv := mapper(key, value)
			return mk, mv, valid
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
		return *new(V), false
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
		return *new(MK), *new(MV), false
	}
}
