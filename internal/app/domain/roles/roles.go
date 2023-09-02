package roles

const (
	RolePUMA       = 1
	RoleComputizen = 2
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
