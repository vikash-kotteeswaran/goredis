package core

var (
	SET_CMD      Command = Command{Command: SET, Executor: setKeyExec}
	GET_CMD      Command = Command{Command: GET, Executor: getKeyExec}
	ECHO_CMD     Command = Command{Command: ECHO, Executor: echoExec}
	PING_CMD     Command = Command{Command: PING, Executor: pingExec}
	REPLCONF_CMD Command = Command{Command: REPLCONF, Executor: replConfExec}
	PSYNC_CMD    Command = Command{Command: PSYNC, Executor: pSyncExec}
	INFO_CMD     Command = Command{Command: INFO, Executor: infoExec}
	ABORT_CMD    Command = Command{Command: ABORT, Executor: abortExec}
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
