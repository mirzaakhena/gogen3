package command2

type Runner interface {
	Run(inputs ...string) error
}
