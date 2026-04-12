package core

const VERSION = "0.2.2"

type Status int

const (
	SUCCESS Status = 1
	FAILURE Status = 2
	RUNNING Status = 3
	ERROR   Status = 4
)

const (
	COMPOSITE = "composite"
	DECORATOR = "decorator"
	ACTION    = "action"
	CONDITION = "condition"
)
