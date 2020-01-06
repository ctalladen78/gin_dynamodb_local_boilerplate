package main

type TodoObject struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdat"`
	Todo      string `json:"todo"`
}

type UserObject struct {
	Name string `json:"name"`
}
