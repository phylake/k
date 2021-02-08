package cmd

var cmds map[string]func()

func Register(name string, f func()) {
	if cmds == nil {
		cmds = make(map[string]func())
	}
	cmds[name] = f
}

func GetCommand(name string) (f func(), exists bool) {
	f, exists = cmds[name]
	return
}
