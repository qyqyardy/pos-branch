package model

import "time"

type StoreSettings struct {
	Name         string
	Tagline      string
	AddressLine1 string
	AddressLine2 string
	Phone        string
	LogoDataURL  string
	Plan         string
	PaidUntil    time.Time
}
