package pipelinerunner

type PipelineRun struct {
	ID          string    `db:"id"`
	Pipelineref string    `db:"pipelineref"`
	Start       string    `db:"start"`
	End         string    `db:"end"`
	Status      RunStatus `db:"status"`
}

func NewPipelineRun(pipelineref string, start string, end string, status RunStatus) *PipelineRun {
	pr := new(PipelineRun)
	pr.ID = getUUID()
	pr.Pipelineref = pipelineref
	pr.Start = start
	pr.End = end
	pr.Status = status
	return pr
}

func (pr PipelineRun) Equal(other PipelineRun) bool {
	return pr.Pipelineref == other.Pipelineref &&
		pr.Start == other.Start &&
		pr.End == other.End &&
		pr.Status == other.Status
}
