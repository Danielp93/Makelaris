package collectors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type EngelEnVolkers struct {
	name           string
	connectionInfo *ConnectionInfo
}

func NewEngelEnVolkers() *EngelEnVolkers {
	return &EngelEnVolkers{
		name: "EngelEnVolkers Amsterdam",
		connectionInfo: &ConnectionInfo{
			RentURL: "https://www.engelvoelkers.com/en/search/?q=&startIndex=0&businessArea=residential&sortOrder=DESC&sortField=sortPrice&pageSize=18&facets=dstrct%3Aamsterdam%3Brgn%3Anoord_holland%3Bcntry%3Anetherlands%3Bbsnssr%3Aresidential%3Bobjcttyp%3Acondo%3Btyp%3Arent%3B",
		},
	}
}

func (d *EngelEnVolkers) Name() string {
	return d.name
}

func (d *EngelEnVolkers) ConnInfo() *ConnectionInfo {
	return d.connectionInfo
}

func (d *EngelEnVolkers) RentingFunction() (string, func(s *colly.HTMLElement)) {
	return ".ev-property-container", func(s *colly.HTMLElement) {
		if strings.Contains(s.Attr("class"), "ev-newsletter-widget") || strings.Contains(s.Attr("class"), "ev-banner-widget") {
			return
		}
		fmt.Print(s.ChildText(".ev-teaser-title"))
		fmt.Print(": ")
		fmt.Println(s.Attr("href"))
	}
}
