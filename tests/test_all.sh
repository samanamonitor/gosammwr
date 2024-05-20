#!/bin/bash

test_run_pipe() {
	NAME=$1
	shift
	INPUT=$1
	shift
	OUTPUT=${NAME}.txt
	cat ${INPUT} | $@ > ${OUTPUT}
	res=$?
	if [ "${res}" != "0" ]; then
		echo "Test ${NAME} failed - Execution Failed." >&2
		exit 1
	fi
	echo "Test ${NAME} successful."
}

test_output() {
	NAME=$1
	shift
	OUTPUT=$1
	shift
	VALID=$1
	diff $OUTPUT $VALID
	res=$?
	if [ "${res}" != "0" ]; then
		echo "Test ${NAME} failed - Output was different." >&2
		exit 1
	fi
	echo "Test ${NAME} successful."

}

test_run_pipe get1 /dev/null bin/test_protocol_get http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/* \
	__cimnamespace=root/cimv2 ClassName=Win32_OperatingSystem
test_output get1_output get1.txt valid/valid1.txt
rm get1.txt

test_run_pipe get2 /dev/null bin/test_protocol_get http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/win32_diskdrive \
	"DeviceId=\\\\.\\PHYSICALDRIVE2"
test_output get2_output get2.txt valid/valid2.txt
rm get2.txt

test_run_pipe create /dev/null bin/test_protocol_create http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd
test_run_pipe command create.txt bin/test_protocol_command - echo ping
test_run_pipe receive command.txt bin/test_protocol_receive -
test_run_pipe delete create.txt bin/test_protocol_delete -
test_output receive_output receive.txt valid/valid3.txt
rm create.txt command.txt receive.txt delete.txt

test_run_pipe enum_computersystem /dev/null bin/test_protocol_enum resourceuri http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/win32_ComputerSystem
test_run_pipe pull_computersystem enum_computersystem.txt bin/test_protocol_pull -
test_output pull_computersystem_output pull_computersystem.txt valid/valid4.txt
rm pull_computersystem.txt enum_computersystem.txt

test_run_pipe enum_diskdrive /dev/null bin/test_protocol_enum resourceuri http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/win32_DiskDrive \
	index 1
test_run_pipe pull_diskdrive enum_diskdrive.txt bin/test_protocol_pull -
test_output pull_diskdrive_output pull_diskdrive.txt valid/valid5.txt
rm enum_diskdrive.txt pull_diskdrive.txt

#test_run_pipe enum_schema /dev/null bin/test_protocol_enum schema root/cimv2 win32_ComputerSystem
#test_run_pipe pull_schema enum_schema.txt bin/test_protocol_pull -
#test_output pull_schema_output pull_schema.txt valid/valid6.txt
#rm enum_schema.txt pull_schema.txt