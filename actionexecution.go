package pipelinerunner

type ActionExecution struct {
	ID             string    `db:"id"`
	Actionref      string    `db:"actionref"`
	Pipelinerunref string    `db:"pipelinerunref"`
	Start          string    `db:"start"`
	End            string    `db:"end"`
	Stdout         string    `db:"stdout"`
	Stderr         string    `db:"stderr"`
	Status         RunStatus `db:"status"`
}

func NewActionExecution(
	actionref string,
	pipelinerunref string,
	start string,
	end string,
	stdout string,
	stderr string,
	status RunStatus,
) *ActionExecution {
	ae := new(ActionExecution)
	ae.ID = getUUID()
	ae.Actionref = actionref
	ae.Pipelinerunref = pipelinerunref
	ae.Start = start
	ae.End = end
	ae.Stdout = stdout
	ae.Stderr = stderr
	ae.Status = status
	return ae
}

func (ae ActionExecution) Equal(other ActionExecution) bool {
	return ae.Actionref == other.Actionref &&
		ae.Pipelinerunref == other.Pipelinerunref &&
		ae.Start == other.Start &&
		ae.End == other.End &&
		ae.Stdout == other.Stdout &&
		ae.Stderr == other.Stderr &&
		ae.Status == other.Status
}
