package core

import "database/sql"

type Provider struct {
	id   int
	name string
}

///  Getters

func (p *Provider) Id() int {
	return p.id
}

func (p *Provider) Name() string {
	return p.name
}

/// Setters

func (p *Provider) SetId(id int) {
	p.id = id
}

func (p *Provider) SetName(name string) {
	p.name = name
}

// [Scan] return object of [Sim] whitch is [Scannable], and map index [int]
// If any errors ocured while scanning it will be in [error]
func (p *Provider) ScanRows(row *sql.Rows) (int, error) {
	err := row.Scan(&p.id, &p.name)
	return p.id, err
}

func (p *Provider) ScanRow(row *sql.Row) error {
	err := row.Scan(&p.id, &p.name)
	return err
}

func (p *Provider) GetKey() int {
	return p.Id()
}

func (p *Provider) SetKey(id int) {
	p.id = id
}
