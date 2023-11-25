package class

type Class struct {
	val string
}

var (
	Unknown = Class{}
	SS4     = Class{val: "ss-4"}
	T4      = Class{val: "t-4"}
	NS4     = Class{val: "ns-4"}
	S4      = Class{val: "ns-4"}
)

func (c Class) String() string {
	return c.val
}
