package rest

type hasDomain[T any] interface {
	toDomain() T
}

func toDomainPointer[T any, D hasDomain[T]](domain *D) *T {
	if domain == nil {
		return nil
	}
	value := (*domain).toDomain()
	return &value
}

func toDomainSlice[T any, D hasDomain[T]](domains []D) []T {
	slice := make([]T, len(domains))
	for i := range domains {
		slice[i] = domains[i].toDomain()
	}
	return slice
}

func toDomainSlicePointer[T any, D hasDomain[T]](domains *[]D) *[]T {
	if domains == nil {
		return nil
	}
	actual := *domains
	slice := make([]T, len(actual))
	for i := range actual {
		slice[i] = actual[i].toDomain()
	}
	return &slice
}
