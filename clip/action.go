package clip 

type Action int

const (
	APPEND Action = iota +1
	SELECT
	DELETE
)

func (a Action) String() string {
	return [...]string{"APPEND", "SELECT", "DELETE"}[a-1]
}

func (a Action) EnumIndex() int{
	return int(a)
}