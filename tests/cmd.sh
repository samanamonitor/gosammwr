#!/bin/bash

set -x

SHELL_ID=$(./test_shell_create)
if [ "$?" != "0" ]; then
	exit 1
fi

COMMAND_ID=$(WINRS_SKIP_CMD_SHELL=1 ./test_shell_command $SHELL_ID $@)
if [ "$?" != "0" ]; then
	./test_shell_delete ${SHELL_ID}
fi

end=0
while [ $end == 0 ]; do
	./test_shell_receive ${SHELL_ID} ${COMMAND_ID} stdout
	if [ $? == 0 ]; then
		end=1
	fi
done

./test_shell_delete ${SHELL_ID}
