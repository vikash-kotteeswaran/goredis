package core

type CommandType struct {
	CREAD     CType
	CWRITE    CType
	CINFO     CType
	CREPL     CType
	CCRITICAL CType
}

var CTYPE = CommandType{CREAD: 1, CWRITE: 2, CINFO: 4, CREPL: 8, CCRITICAL: 16}

var (
	SET_CMD      Command = Command{Command: SET, Type: CTYPE.CWRITE, Executor: setKeyExec}
	GET_CMD      Command = Command{Command: GET, Type: CTYPE.CREAD, Executor: getKeyExec}
	ECHO_CMD     Command = Command{Command: ECHO, Type: CTYPE.CINFO, Executor: echoExec}
	PING_CMD     Command = Command{Command: PING, Type: CTYPE.CINFO, Executor: pingExec}
	REPLCONF_CMD Command = Command{Command: REPLCONF, Type: CTYPE.CREPL, Executor: replConfExec}
	PSYNC_CMD    Command = Command{Command: PSYNC, Type: CTYPE.CREPL, Executor: pSyncExec}
	INFO_CMD     Command = Command{Command: INFO, Type: CTYPE.CINFO, Executor: infoExec}
	ABORT_CMD    Command = Command{Command: ABORT, Type: CTYPE.CCRITICAL, Executor: abortExec}
)

var CommandMap map[string]Command = map[string]Command{
	SET:      SET_CMD,
	GET:      GET_CMD,
	ECHO:     ECHO_CMD,
	PING:     PING_CMD,
	REPLCONF: REPLCONF_CMD,
	PSYNC:    PSYNC_CMD,
	INFO:     INFO_CMD,
	ABORT:    ABORT_CMD,
}
