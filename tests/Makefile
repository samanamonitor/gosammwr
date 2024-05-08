BINFILES=test_protocol_enum \
	test_protocol_enum_cmd \
	test_protocol_get \
	test_protocol_release \
	test_protocol_pull \
	test_protocol_create \
	test_protocol_delete \
	test_protocol_receive \
	test_protocol_send \
	test_mime_encode \
	test_shell_list \
	test_shell_create \
	test_shell_delete

all: $(BINFILES)

%: %.go
	go build $<

test: test.go
	go build $< 

clean:
	rm -f $(BINFILES)
