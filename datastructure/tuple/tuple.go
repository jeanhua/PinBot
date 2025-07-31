package tuple

type Tuple[T comparable, U comparable] struct {
	First  T
	Second U
}

func Of[T comparable, U comparable](first T, second U) Tuple[T, U] {
	return Tuple[T, U]{
		First:  first,
		Second: second,
	}
}

func (t *Tuple[T, U]) Equal(other *Tuple[T, U]) bool {
	return t.First == other.First && t.Second == other.Second
}
