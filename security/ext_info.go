package security

type ExtInfo struct {
	Module string
	Perms  []string
}

func (ext *ExtInfo) Ext(perms ...string) *ExtInfo {
	if len(perms) == 0 {
		return ext
	} else {
		return &ExtInfo{Module: ext.Module, Perms: perms}
	}
}
