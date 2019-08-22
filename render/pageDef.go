package render

import (
	"github.com/kllpw/kllpw-web/ascii"
)
var header = ascii.RenderString(" kll.pw")

var Index = Page{
	Filename:     "index.html",
	Title:        "kllpw",
	Header:       header,
	ContentTitle: "",
	Content:      nil,
}

var Login = Page{
	Filename:     "login.html",
	Title:        "kllpw",
	Header:       header,
	ContentTitle: "",
	Content:      nil,
}

var Register = Page{
	Filename:     "register.html",
	Title:        "kllpw",
	Header:       header,
	ContentTitle: "",
	Content:      nil,
}

var UserHome = Page{
	Filename:     "userHome.html",
	Title:        "kllpw",
	Header:       "",
	ContentTitle: "",
	Content:      nil,
}
