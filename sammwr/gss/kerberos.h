#ifndef _GOKERBEROS
#define _GOKERBEROS
#include <gssapi/gssapi.h>
#include "kerberos.h"

#define AUTH_GSS_ERROR      -1
#define AUTH_GSS_COMPLETE    1
#define AUTH_GSS_CONTINUE    0

#define GSS_AUTH_P_NONE         1
#define GSS_AUTH_P_INTEGRITY    2
#define GSS_AUTH_P_PRIVACY      4

typedef struct {
    gss_ctx_id_t     context;
    gss_name_t       server_name;
    gss_OID          mech_oid;
    long int         gss_flags;
    gss_cred_id_t    client_creds;
    char*            username;
    char*            response;
    size_t           response_len;
    int              responseConf;
    uint32_t         maj_stat;
    uint32_t         min_stat;
    const char *     err_msg;
} gss_client_state;

typedef struct {
    gss_ctx_id_t     context;
    gss_name_t       server_name;
    gss_name_t       client_name;
    gss_cred_id_t    server_creds;
    gss_cred_id_t    client_creds;
    char*            username;
    char*            targetname;
    char*            response;
    size_t           response_len;
    uint32_t         maj_stat;
    uint32_t         min_stat;
    const char *     err_msg;
} gss_server_state;

const char* serverPrincipalDetails(const char *service, const char *hostname);

const char * auth_gss_error_get(void *state, uint32_t *maj_stat, uint32_t *min_stat);
void * auth_gss_new_client();
int auth_gss_client_init(void *s,
        const char *service,
        const char *principal, 
        const char *password, 
        const char *keytab_file, 
        uint32_t gss_flags);
int auth_ggs_client_clean(void *s);
int auth_gss_client_step(void *s, char* challenge, size_t challenge_len, struct gss_channel_bindings_struct* channel_bindings);
char * auth_gss_client_response(void *s, size_t *response_len);
int auth_gss_client_wrap_iov(void *s, char *data, size_t data_len, int protect, int *pad_len);
int auth_gss_client_unwrap_iov(void *s, char *data, size_t len);
void *channel_bindings(OM_uint32 initiator_addrtype,
        size_t initiator_address_length,
        void *initiator_address_value,
        OM_uint32 acceptor_addrtype,
        size_t  acceptor_address_length,
        void *acceptor_address_value,
        size_t application_data_length,
        void *application_data_value);
void * auth_gss_new_server();
int auth_gss_server_init(void *s, char *service,
        const char *password, 
        const char *keytab_file);
int auth_gss_server_clean(void *s);
void debug_gss_client_state(void *s);

#endif /* _GOKERBEROS */