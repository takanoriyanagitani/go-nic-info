package itools

import (
	"iter"
)

func ToSeq2[K,V any](
	original iter.Seq[K],
	mapper func(K) V,
) iter.Seq2[K,V] {
	return func(yield func(K,V) bool){
		for key := range original {
			var val V = mapper(key)
			if !yield(key, val){
				return
			}
		}
	}
}
