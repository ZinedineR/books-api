package externalapi

type ExampleSvcExternal interface {
	Post() (interface{}, int, error)
}
