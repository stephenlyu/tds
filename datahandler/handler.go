package datahandler

import "github.com/stephenlyu/tds/entity"

type RecordHandler interface {
	Feed(r *entity.Record) []*entity.Record
}
