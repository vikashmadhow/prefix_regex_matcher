package seq

import "iter"

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

func FilterSeq[V any](seq iter.Seq[V], predicate FilterFunc[V]) iter.Seq[V] {
	return func(yield func(value V) bool) {
		for v := range seq {
			if predicate(v) {
				if !yield(v) {
					break
				}
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

func FilterSeq2[K, V any](seq iter.Seq2[K, V], predicate Filter2Func[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if predicate(k, v) {
				if !yield(k, v) {
					break
				}
			}
		}
	}
}
