package rest

import (
	"reflect"
)

type ReflectConsumer interface {
	Consume(v reflect.Value)
	ForIndex(index int) ReflectConsumer
	ForField(field reflect.StructField) ReflectConsumer
	ForKey(key string) ReflectConsumer
}

type ReflectShouldConsume func(t reflect.Type) bool

type ReflectIterator func(value reflect.Value, consumer ReflectConsumer)

type Reflector struct {
	Iterators     map[reflect.Type][]ReflectIterator
	ShouldConsume ReflectShouldConsume
}

func NewReflector(should ReflectShouldConsume) *Reflector {
	return &Reflector{
		Iterators:     make(map[reflect.Type][]ReflectIterator),
		ShouldConsume: should,
	}
}

func (r *Reflector) Consume(value any, consumer ReflectConsumer) {
	reflectValue := reflect.ValueOf(value)
	reflectType := reflectValue.Type()
	iterators := r.getIterators(reflectType)
	r.iterate(iterators, reflectValue, consumer)
}

func (r *Reflector) getIterators(t reflect.Type) []ReflectIterator {
	iterators, exists := r.Iterators[t]

	if exists {
		return iterators
	}

	r.Iterators[t] = make([]ReflectIterator, 0)

	if r.ShouldConsume(t) {
		r.addIterator(t, r.consumeValue())
	}

	switch t.Kind() {
	case reflect.Pointer:
		r.addIterator(t, r.consumePointer(t.Elem()))
	case reflect.Interface:
		r.addIterator(t, r.consumeInterface(t.Elem()))
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			r.addIterator(t, r.consumeField(t, i))
		}
	case reflect.Array, reflect.Slice:
		r.addIterator(t, r.consumeSlice(t.Elem()))
	case reflect.Map:
		r.addIterator(t, r.consumeMap(t.Elem()))
	}

	return r.Iterators[t]
}

func (r *Reflector) addIterator(t reflect.Type, iter ReflectIterator) {
	if iter != nil {
		r.Iterators[t] = append(r.Iterators[t], iter)
	}
}

func (r *Reflector) iterate(iters []ReflectIterator, value reflect.Value, consumer ReflectConsumer) {
	for _, i := range iters {
		i(value, consumer)
	}
}

func (r *Reflector) consumeValue() ReflectIterator {
	return func(value reflect.Value, consumer ReflectConsumer) {
		consumer.Consume(value)
	}
}

func (r *Reflector) consumePointer(t reflect.Type) ReflectIterator {
	iters := r.getIterators(t)
	if len(iters) == 0 {
		return nil
	}

	return func(value reflect.Value, consumer ReflectConsumer) {
		if value.IsNil() {
			return
		}
		r.iterate(iters, value.Elem(), consumer)
	}
}

func (r *Reflector) consumeInterface(t reflect.Type) ReflectIterator {
	iters := r.getIterators(t)
	if len(iters) == 0 {
		return nil
	}

	return func(value reflect.Value, consumer ReflectConsumer) {
		if value.IsNil() {
			return
		}
		r.iterate(iters, value.Elem(), consumer)
	}
}

func (r *Reflector) consumeSlice(t reflect.Type) ReflectIterator {
	iters := r.getIterators(t)
	if len(iters) == 0 {
		return nil
	}

	return func(value reflect.Value, consumer ReflectConsumer) {
		if value.IsNil() {
			return
		}
		for i := 0; i < value.Len(); i++ {
			r.iterate(iters, value.Index(i), consumer.ForIndex(i))
		}
	}
}

func (r *Reflector) consumeMap(t reflect.Type) ReflectIterator {
	iters := r.getIterators(t)
	if len(iters) == 0 {
		return nil
	}

	return func(value reflect.Value, consumer ReflectConsumer) {
		if value.IsNil() {
			return
		}
		iter := value.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()

			r.iterate(iters, value, consumer.ForKey(key.String()))
		}
	}
}

func (r *Reflector) consumeField(parent reflect.Type, fieldIndex int) ReflectIterator {
	field := parent.Field(fieldIndex)
	iters := r.getIterators(field.Type)

	if len(iters) == 0 {
		return nil
	}

	return func(value reflect.Value, consumer ReflectConsumer) {
		fieldValue := value.Field(fieldIndex)
		if fieldValue.CanAddr() {
			fieldValue = fieldValue.Addr()
		}
		if field.Anonymous {
			r.iterate(iters, fieldValue, consumer)
		} else {
			r.iterate(iters, fieldValue, consumer.ForField(field))
		}
	}
}
