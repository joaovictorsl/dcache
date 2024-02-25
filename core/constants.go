package core

const (
	CMD_SET byte = iota
	CMD_GET
	CMD_HAS
	CMD_DELETE

	CMD_EXEC_SUCCEEDED
	CMD_EXEC_FAILED

	INVALID_COMMAND_CODE
	INVALID_COMMAND string = "invalid command"
)
