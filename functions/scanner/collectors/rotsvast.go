package collectors

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type RotsVast struct {
	name           string
	connectionInfo *ConnectionInfo

	collection []*Listing
}

func NewRotsVast() *RotsVast {
	return &RotsVast{
		name: "Rotsvast Amsterdam",
		connectionInfo: &ConnectionInfo{
			RentURL: "https://www.rotsvast.nl/woningaanbod/rotsvast-amsterdam/?type=2&&display=list&count=1000",
		},
		collection: []*Listing{},
	}
}

func (d *RotsVast) Name() string {
	return d.name
}

func (d *RotsVast) ConnInfo() *ConnectionInfo {
	return d.connectionInfo
}

func (d *RotsVast) Collect(c colly.Collector) []*Listing {
	//Parsing Rental Listings
	c.OnHTML(".residence-list", func(e *colly.HTMLElement) {

		//Extract streetname and house
		street := e.ChildText(".residence-street")

		//Extract Zipcode and city
		fullCity := e.ChildText(".residence-zipcode-place")
		zipcode := fullCity[0:6]
		city := fullCity[7:]

		//Extract Price Variable
		splitPrice := strings.Split(e.ChildText(".residence-price"), ",")
		p := strings.TrimPrefix(splitPrice[0], "€ ")
		p = strings.ReplaceAll(p, ".", "")
		price, _ := strconv.Atoi(p)

		//Extract Size
		p = strings.TrimSuffix(e.ChildText(".residence-properties .row .col-md-6"), " m")
		p = strings.TrimPrefix(p, "Woonoppervlakte")
		size, _ := strconv.Atoi(p)

		listing := &Listing{
			Address: &ListingAddress{
				Street: street,
				City:   city,
				Zip:    zipcode,
			},
			Type:   "Huur",
			status: e.ChildText(".status"),
			Price:  price,
			Size:   size,
			Url:    e.ChildAttr(".clickable-block", "href"),
		}
		d.collection = append(d.collection, listing)
	})
	c.Visit(d.connectionInfo.RentURL)
	return d.collection
}
