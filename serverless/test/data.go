package test

func ConfigureDb() {
	AddToEnv(DbEnv)
}

func ExampleDbRecords() [][]string {
	return [][]string{
		{"teams", "Fox", "", "", "", "", "GPAVideo"},
		{"admins", "Fox", "admin@example.com", "John", "Public", "admin", ""},
		{"guests", "Fox", "guest@example.com", "Kristy", "Thomas", "guest", ""},
	}
}
