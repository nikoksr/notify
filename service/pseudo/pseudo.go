package pseudo

type Pseudo struct{}

func New() *Pseudo {
	return &Pseudo{}
}

func (Pseudo) Send(string, string) error {
	return nil
}
