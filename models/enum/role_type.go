package enum

type RoleType int8

const (
	AdminRole RoleType = iota + 1
	UserRole
)
