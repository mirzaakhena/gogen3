package driver

type Controller interface {
	RegisterRouter()
}

type RegistryContract interface {
	Controller
	RunApplication()
}

func Run(rv RegistryContract) {
	if rv != nil {
		rv.RegisterRouter()
		rv.RunApplication()
	}
}
