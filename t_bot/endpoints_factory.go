package t_bot

type Endpoints interface {
	GetCrime(idParam string)
}

func NewEndpointsFactory(crimeEvent CrimeEvents) *endpointsFactory {
	return &endpointsFactory{crimeEvents:crimeEvent}
}

type endpointsFactory struct {
	crimeEvents CrimeEvents
}
