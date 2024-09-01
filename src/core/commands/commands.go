package commands

var (
	SET_CMD   Command = Command{command: SET, executor: setKeyExec}
	GET_CMD   Command = Command{command: GET, executor: getKeyExec}
	ECHO_CMD  Command = Command{command: ECHO, executor: echoExec}
	PING_CMD  Command = Command{command: PING, executor: pingExec}
	INFO_CMD  Command = Command{command: INFO, executor: infoExec}
	ABORT_CMD Command = Command{command: ABORT, executor: abortExec}
)

var CommandMap map[string]Command = map[string]Command{
	SET:   SET_CMD,
	GET:   GET_CMD,
	ECHO:  ECHO_CMD,
	PING:  PING_CMD,
	INFO:  INFO_CMD,
	ABORT: ABORT_CMD,
}
