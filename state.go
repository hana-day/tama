package tama

type State struct {
}

func NewState() *State {
	return &State{}
}

func (s *State) ExecString(source string) (result string, err error) {
	return source, nil
}
