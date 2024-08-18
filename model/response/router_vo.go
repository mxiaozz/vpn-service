package response

type RouterVo struct {
	Name       string     `json:"name,omitempty"`
	Path       string     `json:"path,omitempty"`
	Hidden     bool       `json:"hidden,omitempty"`
	Redirect   string     `json:"redirect,omitempty"`
	Component  string     `json:"component,omitempty"`
	Query      string     `json:"query,omitempty"`
	AlwaysShow bool       `json:"alwaysShow,omitempty"`
	Meta       *MetaVo    `json:"meta,omitempty"`
	Children   []RouterVo `json:"children,omitempty"`
}
