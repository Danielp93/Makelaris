package main

import "scanner/collectors"

func main() {
	c := collectors.New()

	c.AddAgents(
		// collectors.NewDeGraafEnGroot(),
		// collectors.NewEngelEnVolkers(),
		collectors.NewRotsVast(),
		//\collectors.NewSotherby(),
	)

	c.Collect()
}
