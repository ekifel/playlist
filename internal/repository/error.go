package repository

type ObjNotFound struct {
	Msg string
}

func (e *ObjNotFound) Error() string {
	return e.Msg
}
