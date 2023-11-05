package listing

//go:generate stringer -type Status

type Status int

const (
	beschikbaar Status = iota
	verhuurd
	verkocht
	geboden
	uitgelicht
)
