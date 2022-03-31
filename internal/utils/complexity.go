package utils

import "github.com/saime-0/messenger-for-employee/graph/generated"

func add(v int) func(childComplexity int) int {
	return func(childComplexity int) int {
		return childComplexity + v
	}
}
func mul(v int) func(childComplexity int) int {
	return func(childComplexity int) int {
		return childComplexity * v
	}
}

func MatchComplexity() *generated.ComplexityRoot {
	c := &generated.ComplexityRoot{}
	c.Room.Members = add(4)
	//c.Room.Messages = add(5)
	c.Messages.Messages = add(5)
	c.Members.Members = add(4)
	c.Rooms.Rooms = add(4)
	c.Employees.Employees = add(3)
	return c
}
