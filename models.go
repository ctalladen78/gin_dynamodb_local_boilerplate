package main

type TodoObject struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdat"`
	Todo      string `json:"todo"`
}

type UserObject struct {
	Name string `json:"name"`
}

// https://yourbasic.org/golang/iota/
// enum for querying item attributes
type QueryCondition int

const (
	CREATED_AT = iota
	CREATED_BY
)

func (q QueryCondition) String() string {
	return [...]string{"CREATED_BY", "CREATED_AT"}[q]
}
