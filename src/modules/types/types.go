package types

import "time"

type Snapshot struct {
	ID       string   `json:"id"`
	Time     string   `json:"time"`
	Tree     string   `json:"tree"`
	Paths    []string `json:"paths"`
	Hostname string   `json:"hostname"`
	Username string   `json:"username"`
	UID      int      `json:"uid"`
	GID      int      `json:"gid"`
	ShortID  string   `json:"short_id"`
}

type File struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Path        string    `json:"path"`
	UID         int       `json:"uid"`
	GID         int       `json:"gid"`
	Mode        int       `json:"mode"`
	Permissions string    `json:"permissions"`
	MTime       time.Time `json:"mtime"`
	ATime       time.Time `json:"atime"`
	CTime       time.Time `json:"ctime"`
	StructType  string    `json:"struct_type"`
}
