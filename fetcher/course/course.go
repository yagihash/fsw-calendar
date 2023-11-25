package course

type Course struct {
	val string
}

var (
	Unknown = Course{}
	RC      = Course{val: "rc"} // Racing course
	SS      = Course{val: "ss"} // Short course
)

func (c Course) String() string {
	return c.val
}
