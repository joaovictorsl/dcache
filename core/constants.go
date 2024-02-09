package core

const (
	CMD_SET byte = iota
	CMD_GET
	CMD_HAS
	CMD_DELETE

	CMD_EXEC_SUCCEEDED
	CMD_EXEC_FAILED
	INVALID_COMMAND_CODE

	INVALID_COMMAND        string = "invalid command"
	INVALID_GET_COMMAND    string = "invalid GET command"
	INVALID_SET_COMMAND    string = "invalid SET command"
	INVALID_HAS_COMMAND    string = "invalid HAS command"
	INVALID_DELETE_COMMAND string = "invalid DELETE command"
)
