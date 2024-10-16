package types

func DefaultParams() Params {
	return Params{
		EnableErc20: true,
	}
}

func (p *Params) Validate() error {
	return nil
}
