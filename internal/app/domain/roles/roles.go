package roles

const (
	RolePUMA = "PUMA"
	RoleUser = "Computizen"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
