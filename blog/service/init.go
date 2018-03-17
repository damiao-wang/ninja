package service

import (
	"ninja/blog/service/article"

	"ninja/base/misc/cmdt"
)

func Init() {
	cmdt.RegisterService(&article.Service{})
}
