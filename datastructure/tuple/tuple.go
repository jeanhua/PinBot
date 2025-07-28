package tuple

type Tuple[FirstValue interface{}, SecondValue interface{}] struct {
	First  FirstValue
	Second SecondValue
}

func Of[FirstValue interface{}, SecondValue interface{}](first FirstValue, second SecondValue) Tuple[FirstValue, SecondValue] {
	return Tuple[FirstValue, SecondValue]{
		First:  first,
		Second: second,
	}
}
