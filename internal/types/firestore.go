package types

import "time"

type Hostname struct {
	Verified bool `firestore:"verified" json:"verified"`
}

type EdgeLogic struct {
	RedirectTo   string    `firestore:"redirect_to" json:"redirect_to"`
	EnforceHTTPS string    `firestore:"enforce_https" json:"enforce_https"`
	CacheControl string    `firestore:"cache_control" json:"cache_control"`
	Created      time.Time `firestore:"created" json:"created"`
	Updated      time.Time `firestore:"updated" json:"updated"`
	Backend      string    `firestore:"backend" json:"backend"`
	BuildID      string    `firestore:"build_id" json:"build_id"`
	Jurisdiction string    `firestore:"jurisdiction" json:"jurisdiction"`
}

type HostnameMetadata struct {
	Hostname string    `firestore:"hostname" json:"hostname"`
	Zone     string    `firestore:"zone" json:"zone"`
	Created  time.Time `firestore:"created" json:"created"`
	Updated  time.Time `firestore:"updated" json:"updated"`
	SiteID   string    `firestore:"site_id" json:"site_id"`
	SiteEnv  string    `firestore:"site_env" json:"site_env"`
}
