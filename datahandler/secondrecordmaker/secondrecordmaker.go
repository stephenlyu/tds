package secondrecordmaker

import (
	"github.com/stephenlyu/tds/datahandler/secondrecordgenerator"
	"github.com/stephenlyu/tds/datahandler/secondsmoother"
	"github.com/stephenlyu/tds/datahandler"
	"github.com/stephenlyu/tds/entity"
)

type SecondRecordMaker struct {
	generator *secondrecordgenerator.SecondRecordGenerator
	smoother datahandler.RecordHandler
}

func NewSecondRecordMaker(security *entity.Security) *SecondRecordMaker {
	return &SecondRecordMaker{
		generator: secondrecordgenerator.NewSecondRecordGenerator(security),
		smoother: secondsmoother.NewSecondSmoother(security, nil),
	}
}

func (this *SecondRecordMaker) Feed(tick *entity.TickItem) []*entity.Record {
	r := this.generator.Feed(tick)
	if r == nil {
		return []*entity.Record{}
	}
	return this.smoother.Feed(r)
}
