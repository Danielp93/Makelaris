package collectors

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type DeGraafEnGroot struct {
	name           string
	connectionInfo *ConnectionInfo
	collection     []*Listing
}

func NewDeGraafEnGroot() *DeGraafEnGroot {
	return &DeGraafEnGroot{
		name: "De Graaf en Groot",
		connectionInfo: &ConnectionInfo{
			RentURL: "https://www.degraafengroot.nl/rentals",
		},
		collection: []*Listing{},
	}
}

func (d *DeGraafEnGroot) Name() string {
	return d.name
}

func (d *DeGraafEnGroot) ConnInfo() *ConnectionInfo {
	return d.connectionInfo
}

func (d *DeGraafEnGroot) Collect(c colly.Collector) []*Listing {
	//Parsing Rental Listings
	c.OnHTML(".teaser--property", func(e *colly.HTMLElement) {

		//Extract streetname and house
		nrRegex := regexp.MustCompile(`\s(\d.*)$`)
		cityAddress := strings.Split(e.ChildText("teaser__title"), " - ")
		city := cityAddress[0]
		fullAddress := cityAddress[1]
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
				City:   city,
				Zip:    "",
			},
			Type:   "Huur",
			status: e.ChildText(".item-content ul li:nth-child(2) span:nth-child(2)"),
			Price:  price,
			Size:   size,
			Url:    e.ChildAttr(".more-information", "href"),
		}
		d.collection = append(d.collection, listing)
	})
	c.Visit(d.connectionInfo.RentURL)
	return d.collection
}
