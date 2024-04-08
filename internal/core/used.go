package core

import (
	"database/sql"
)

type Used struct {
	id          int
	simId       int
	serviceId   int
	isBlocked   bool
	blockedInfo string
}

/// Getters

func (u *Used) Id() int {
	return u.id
}
func (u *Used) SimID() int {
	return u.simId
}
func (u *Used) ServiceID() int {
	return u.serviceId
}
func (u *Used) IsBlocked() bool {
	return u.isBlocked
}

func (u *Used) BlockedInfo() string {
	return u.blockedInfo
}

// / With
func (u Used) WithSimID(id int) Used {
	u.simId = id
	return u
}

func (u Used) WithServiceID(id int) Used {
	u.serviceId = id
	return u
}

// / Setters
func (u *Used) SetId(id int) {
	u.id = id
}
func (u *Used) SetSimID(simID int) {
	u.simId = simID
}
func (u *Used) SetServiceID(sid int) {
	u.serviceId = sid
}
func (u *Used) SetIsBlocked(status bool) {
	u.isBlocked = status
}
func (u *Used) SetBlockedInfo(binfo string) {
	u.blockedInfo = binfo
}

// [Scan] return object of [Sim] whitch is [Scannable], and map index [int]
// If any errors ocured while scanning it will be in [error]
func (u *Used) ScanRows(row *sql.Rows) (int, error) {
	err := row.Scan(&u.id, &u.simId, &u.serviceId, &u.isBlocked, &u.blockedInfo)
	return u.id, err
}

func (u *Used) ScanRow(row *sql.Row) error {
	err := row.Scan(&u.id, &u.simId, &u.serviceId, &u.isBlocked, &u.blockedInfo)
	return err
}

func (u *Used) GetKey() int {
	return u.Id()
}

func (u *Used) SetKey(id int) {
	u.id = id
}
