package security

type ExtInfo struct {
	Module string
	Perms  []string
}

func (ext ExtInfo) Ext(perms ...string) ExtInfo {
	ext.Perms = perms
	return ext
}
