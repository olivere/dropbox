package env

var (
	// Defaults holds default values.
	Defaults = struct {
		BulkSize int     // default size for bulk indexing
		MaxTake  int     // maximum number of results to return on pagination
		Company  company // default Company (will be created if no company exists)
		User     user    // default User (will be created if no user exists)
	}{
		BulkSize: 100,
		MaxTake:  500,
		Company: company{
			ID:         1,
			Name:       "Meplato GmbH",
			MPCC:       "meplato",
			MPSC:       "meplato",
			MPBC:       "meplato",
			Country:    "DE",
			VATID:      "DE252280354",
			SAGENumber: "00000",
		},
		User: user{
			ID:        1,
			CompanyID: 1, // same as above
			Type:      "human",
			Name:      "Mr Admin",
			Email:     "admin@meplato.com",
			Password:  "admin",
			Role:      "admin", // change when roles are renamed
			Language:  "de",
		},
	}
)

type company struct {
	ID         int64
	Name       string
	MPCC       string
	MPSC       string
	MPBC       string
	Country    string
	VATID      string
	SAGENumber string
}

type user struct {
	ID        int64
	CompanyID int64
	Type      string
	Name      string
	Email     string
	Password  string
	Role      string
	Language  string
}
