package router

type Operation int

const (
	POST Operation = iota
	GET
	PUT
	PATCH
	DELETE
)

func (o Operation) String() string {
	switch o {
	case POST:
		return "POST"
	case GET:
		return "GET"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	default:
		panic("unknown operation")
	}
}
