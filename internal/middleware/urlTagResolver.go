package middleware

import (
	"os"

	"github.com/Bofry/structproto"
)

var _ structproto.TagResolver = UrlTagResolver

func UrlTagResolver(fieldname, token string) (*structproto.Tag, error) {
	var tag *structproto.Tag
	if token != "-" {
		tag = &structproto.Tag{
			Name: os.ExpandEnv(token),
		}
	}
	return tag, nil
}
