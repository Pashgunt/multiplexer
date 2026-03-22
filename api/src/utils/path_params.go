package apiutils

type PathParam string

func (p PathParam) String() string {
	return string(p)
}

const (
	UUID PathParam = "uuid"
)
