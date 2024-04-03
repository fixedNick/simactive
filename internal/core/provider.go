package core

type Provider struct {
	id   int
	name string
}

func (p *Provider) Id() int {
	return p.id
}

func (p *Provider) Name() string {
	return p.name
}
