package entities

type Repository struct{
	GitURL				string
	Name				string
	User				User
	// CreateAt			Timestamp
}

type User struct{
	Name				string
}

type SearchResult struct{
	Repository			Repository
	FileURL				string
	FileContent			string
}
