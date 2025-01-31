package dependency

type Dependency struct {
	Constructor interface{}
	Interface   interface{}
	Token       string
}
