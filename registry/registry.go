package registry

type Registry struct {
	serviceById map[string]interface{}
}

func New() Registry {
	return Registry{
		serviceById: make(map[string]interface{}),
	}
}

func (self Registry) Set(id string, service interface{}) {
	self.serviceById[id] = service
}

func (self Registry) Get(id string) (service interface{}) {
	return self.serviceById[id]
}
