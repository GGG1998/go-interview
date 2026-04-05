package task

type Task struct {
	Id   string
	Name string
}

func (t Task) GetId() string {
	return t.Id
}
