package services

type Services struct {
	serviceById map[string]interface{}
}

func New() Services {
	return Services{
		serviceById: make(map[string]interface{}),
	}
}

func (self Services) Set(id string, service interface{}) {
	self.serviceById[id] = service
}

func (self Services) Get(id string) (service interface{}) {
	return self.serviceById[id]
}
