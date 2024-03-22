package variable_parameter_test

import "fmt"

func VariableParameterTest() {
	fmt.Println("\033[1;32;40m  \nstart VariableParameterTest---------------------------------- \033[0m")

	x := getMin(1, 3, 2, 0)
	fmt.Printf("The minimum is: %d\n", x)
	slice := []int{7, 9, 3, 5, 1}
	x = getMin(slice...)
	fmt.Printf("The minimum in the slice is: %d\n", x)
}

func getMin(s ...int) int {
	if len(s) == 0 {
		return 0
	}
	min := s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
	}
	return min
}

// //////////////////////////////////////////////////////////////////////////////////
type Person struct {
	id      int
	name    string
	address string
	phone   int
}

func New(id, phone int, name, addr string) Person {
	return Person{
		id:      id,
		name:    name,
		address: addr,
		phone:   phone,
	}
}

type Option func(person *Person)

var defaultPerson = Person{id: -1, name: "-1", address: "-1", phone: -1}

func WithID(id int) Option {
	return func(m *Person) {
		m.id = id
	}
}

func WithName(name string) Option {
	return func(m *Person) {
		m.name = name
	}
}

func WithAddress(addr string) Option {
	return func(m *Person) {
		m.address = addr
	}
}

func WithPhone(phone int) Option {
	return func(m *Person) {
		m.phone = phone
	}
}

func NewByOption(opts ...Option) Person {
	p := defaultPerson
	for _, o := range opts {
		o(&p)
	}
	return p
}

func NewByOptionWithoutID(id int, opts ...Option) Person {
	p := defaultPerson
	p.id = id
	for _, o := range opts {
		o(&p)
	}
	return p
}

func OptionTest() {
	fmt.Println("\033[1;32;40m  \nstart OptionTest---------------------------------- \033[0m")

	p1 := New(1, 123, "person1", "cache1")
	fmt.Println("not use option,message1:", p1)
	p2 := NewByOption(WithID(2), WithName("person2"), WithAddress("cache2"), WithPhone(456))
	fmt.Println("use option,message1:", p2)
	p3 := NewByOptionWithoutID(3, WithAddress("cache3"), WithPhone(789), WithName("person3"))
	fmt.Println("use option without id,message1:", p3)
}
