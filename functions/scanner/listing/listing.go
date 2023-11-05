package listing

type Adress struct {
	Street  string
	Number  string
	Zip     string
	City    string
	Country string
}

type Listing struct {
	Adress Adress
	Status
}
