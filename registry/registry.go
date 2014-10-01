package registry

type Registry struct {
	serviceById map[string]interface{}
}

func New() Registry {
	return Registry{
		serviceById: make(map[string]interface{}),
	}
}

func (self Registry) SetService(id string, service interface{}) {
	self.serviceById[id] = service
}

func (self Registry) GetService(id string) (service interface{}) {
	return self.serviceById[id]
}
