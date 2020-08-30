package formtype

type FormType int

const (
	Form10K FormType = iota
	Form10Q
)

func (d FormType) String() string {
	return [...]string{"10K", "10Q"}[d]
}
