package collectors

import "encoding/json"

type Listing struct {
	Address *ListingAddress
	Type    string
	status  string
	Price   int
	Url     string
	Size    int
}

type ListingAddress struct {
	Street string
	Number string
	City   string
	Zip    string
}

func (l *Listing) prettyPrint() string {
	s, _ := json.MarshalIndent(l, "", "\t")
	return string(s)
}
