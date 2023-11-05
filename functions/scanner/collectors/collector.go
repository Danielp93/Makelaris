package collectors

import (
	"fmt"

	"github.com/gocolly/colly"
)

type ConnectionInfo struct {
	MainURL string
	BuyURL  string
	RentURL string
}

type Agent interface {
	Name() string
	ConnInfo() *ConnectionInfo

	Collect(colly.Collector) []*Listing
}

type Collector struct {
	client *colly.Collector
	agents []Agent

	collection []*Listing
}

func New() *Collector {
	return &Collector{
		client:     colly.NewCollector(),
		agents:     []Agent{},
		collection: []*Listing{},
	}
}

func (c *Collector) AddAgent(a Agent) {
	c.AddAgents(a)
}

func (c *Collector) AddAgents(a ...Agent) {
	for _, agent := range a {
		c.agents = append(c.agents, agent)
	}
}

func (c *Collector) Collect() {
	for _, agent := range c.agents {
		c.collection = append(c.collection, agent.Collect(*c.client)...)
	}

	for _, listing := range c.collection {
		fmt.Println(listing.prettyPrint())
	}
}
