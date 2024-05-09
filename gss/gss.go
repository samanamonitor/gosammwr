package gss

/*
#cgo LDFLAGS: -lkrb5 -lgssapi_krb5
#include "kerberos.h"
*/
import "C"

import (
    "unsafe"
)

const (
    AUTH_GSS_ERROR  = -1
    AUTH_GSS_CONTINUE = 0
    AUTH_GSS_COMPLETE = 1
)

type Gss struct {
    service string
    principal string
    password string
    keytab_file string
    context unsafe.Pointer
    Maj_stat uint
    Min_stat uint
    Error string
}

type ChannelBinding struct {
    data unsafe.Pointer
}

func (self *Gss) AuthGssClientInit(service string, principal string, 
        password string, keytab_file string, gssFlags uint) int {
    self.service = service
    self.principal = principal
    self.password = password
    self.keytab_file = keytab_file        

    self.context = unsafe.Pointer(C.auth_gss_new_client())
    result := int(C.auth_gss_client_init(
            self.context,
            C.CString(self.service), 
            C.CString(self.principal), 
            C.CString(self.password), 
            C.CString(self.keytab_file), 
            C.uint(gssFlags)))
    if result != AUTH_GSS_COMPLETE {
        var maj_stat C.uint
        var min_stat C.uint
        self.Error = C.GoString(C.auth_gss_error_get(self.context, &maj_stat, &min_stat))
        self.Maj_stat = uint(maj_stat)
        self.Min_stat = uint(min_stat)
        return AUTH_GSS_ERROR
    }
    if self.context == nil {
        panic("Context is nil")
    }
    return result
}

func (self *Gss) AuthGssClientClean() int {
    result := int(C.auth_ggs_client_clean(self.context))
    if result != AUTH_GSS_COMPLETE {
        var maj_stat C.uint
        var min_stat C.uint
        self.Error = C.GoString(C.auth_gss_error_get(self.context, &maj_stat, &min_stat))
        self.Maj_stat = uint(maj_stat)
        self.Min_stat = uint(min_stat)
        return AUTH_GSS_ERROR
    }
    return AUTH_GSS_COMPLETE
}

func (self *Gss) AuthGssClientStep(challenge_byte []byte) int {
    if self.context == nil {
        panic("Context is nil")
    }
    var cstr *C.char

    if len(challenge_byte) > 0 {
        cstr = (*C.char)(unsafe.Pointer(&challenge_byte[0]))
    } else {
        cstr = nil
    }
    result := int(C.auth_gss_client_step(self.context, cstr, C.ulong(len(challenge_byte)), nil))
    if result == AUTH_GSS_ERROR {
        var maj_stat C.uint
        var min_stat C.uint
        self.Error = C.GoString(C.auth_gss_error_get(self.context, &maj_stat, &min_stat))
        self.Maj_stat = uint(maj_stat)
        self.Min_stat = uint(min_stat)
    }
    return result
}

func (self *Gss) AuthGssClientResponse() []byte {
    var r *C.char
    var length C.ulong

    r = C.auth_gss_client_response(self.context, &length)
    return C.GoBytes(unsafe.Pointer(r), C.int(length))
}

func (self *Gss) AuthGSSClientWrapIov(message []byte) int {
    var cstr *C.char
    var pad_len C.int

    if len(message) > 0 {
        cstr = (*C.char)(unsafe.Pointer(&message[0]))
    } else {
        cstr = nil
    }

    result := C.auth_gss_client_wrap_iov(self.context, cstr, C.ulong(len(message)), 1, &pad_len)
    if result == AUTH_GSS_ERROR {
        var maj_stat C.uint
        var min_stat C.uint
        self.Error = C.GoString(C.auth_gss_error_get(self.context, &maj_stat, &min_stat))
        self.Maj_stat = uint(maj_stat)
        self.Min_stat = uint(min_stat)
    }
    return int(result)
}

func (self *Gss) AuthGSSClientUnwrapIov(message []byte) int {
    var cstr *C.char

    if len(message) > 0 {
        cstr = (*C.char)(unsafe.Pointer(&message[0]))
    } else {
        cstr = nil
    }

    result := C.auth_gss_client_unwrap_iov(self.context, cstr, C.ulong(len(message)))
    if result == AUTH_GSS_ERROR {
        var maj_stat C.uint
        var min_stat C.uint
        self.Error = C.GoString(C.auth_gss_error_get(self.context, &maj_stat, &min_stat))
        self.Maj_stat = uint(maj_stat)
        self.Min_stat = uint(min_stat)
    }
    return int(result)
}

func (self *Gss) DebugClientState() {
    C.debug_gss_client_state(self.context)
}

func NewChannelBinding(initiator_addrtype uint32, initiator_address []byte, 
        acceptor_addrtype uint32, acceptor_address []byte, 
        application_data []byte) ChannelBinding {

    var initiator_address_ptr unsafe.Pointer
    var acceptor_address_ptr unsafe.Pointer
    var application_data_ptr unsafe.Pointer
    var cb ChannelBinding

    if len(initiator_address) > 0 {
        initiator_address_ptr = unsafe.Pointer(&initiator_address[0])
    } else {
        initiator_address_ptr = nil;
    }
    if len(acceptor_address) > 0 {
        acceptor_address_ptr = unsafe.Pointer(&acceptor_address[0])
    } else {
        acceptor_address_ptr = nil
    }
    if len(application_data) > 0 {
        application_data_ptr = unsafe.Pointer(&application_data[0])
    } else {
        application_data_ptr = nil
    }
    cb.data = C.channel_bindings(C.uint(initiator_addrtype), 
        C.ulong(len(initiator_address)), initiator_address_ptr,
        C.uint(acceptor_addrtype), 
        C.ulong(len(acceptor_address)), acceptor_address_ptr,
        C.ulong(len(application_data)), application_data_ptr)
    return cb
}

