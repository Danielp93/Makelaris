package collectors

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Sotherby struct {
	name           string
	connectionInfo *ConnectionInfo

	collection []*Listing
}

func NewSotherby() *Sotherby {
	return &Sotherby{
		name: "Sotherby International Realty",
		connectionInfo: &ConnectionInfo{
			RentURL: "https://sothebysrealty.nl/aanbod/huur/",
		},
		collection: []*Listing{},
	}
}

func (s *Sotherby) Name() string {
	return s.name
}

func (s *Sotherby) ConnInfo() *ConnectionInfo {
	return s.connectionInfo
}

func (s *Sotherby) Collect(c colly.Collector) []*Listing {
	//Parsing Rental Listings
	c.OnHTML("#entity-items > .woning", func(e *colly.HTMLElement) {

		//Extract streetname and house
		nrRegex := regexp.MustCompile(`\s(\d.*)$`)
		fullAddress := e.ChildText(".item-content > header > h3")
		nrIndex := nrRegex.FindStringSubmatchIndex(fullAddress)
		street := fullAddress[:nrIndex[0]]
		nr := fullAddress[nrIndex[2]:]

		//Extract Price Variable
		p := strings.TrimPrefix(e.ChildText(".item-content ul li:nth-child(1) span:nth-child(2)"), "€ ")
		p = strings.TrimSuffix(p, ",- p/m")
		p = strings.ReplaceAll(p, ".", "")
		price, _ := strconv.Atoi(p)

		//Extract Size
		p = strings.TrimSuffix(e.ChildText(".item-content ul li:nth-child(4) span:nth-child(2)"), " m2")
		size, _ := strconv.Atoi(p)

		listing := &Listing{
			Address: &ListingAddress{
				Street: street,
				Number: nr,
				City:   strings.TrimSuffix(e.ChildText(".item-content > header > h2"), ","),
				Zip:    "",
			},
			Type:   "Huur",
			status: e.ChildText(".item-content ul li:nth-child(2) span:nth-child(2)"),
			Price:  price,
			Size:   size,
			Url:    e.ChildAttr(".more-information", "href"),
		}
		s.collection = append(s.collection, listing)
	})
	c.Visit(s.connectionInfo.RentURL)
	return s.collection
}
