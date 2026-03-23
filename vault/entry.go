package vault

type Entry struct {
	Service   string `json:"service"`
	Account   string `json:"account"`
	Password  string `json:"password"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}
