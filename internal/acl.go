package internal

import (
	"github.com/quzhi1/go-imap/v2"
)

func FormatRights(rm imap.RightModification, rs imap.RightSet) string {
	s := ""
	if rm != imap.RightModificationReplace {
		s = string(rm)
	}
	return s + string(rs)
}
