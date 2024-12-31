package seq

func Filter[V any](seq Seq[V], predicate FilterFunc[V]) Seq[V] {
	return func() (V, bool) {
		for {
			value, valid := seq()
			if !valid {
				return value, valid
			}
			if predicate(value) {
				return value, true
			}
		}
	}
}

func Filter2[K, V any](seq Seq2[K, V], predicate Filter2Func[K, V]) Seq2[K, V] {
	return func() (K, V, bool) {
		for {
			key, value, valid := seq()
			if !valid {
				return key, value, valid
			}
			if predicate(key, value) {
				return key, value, true
			}
		}
	}
}
