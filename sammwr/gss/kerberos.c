#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>
#include <gssapi/gssapi_krb5.h>
#include <gssapi/gssapi_generic.h>
#include "kerberos.h"

const char * auth_gss_error_get(void *s, uint32_t *maj_stat, uint32_t *min_stat)
{
    assert(s != NULL);
    gss_client_state* state = (gss_client_state*) s;
    if (maj_stat != NULL) {
        *maj_stat = state->maj_stat;
    }
    if (min_stat != NULL) {
        *min_stat = state->min_stat;
    }
    return state->err_msg;
}

void * auth_gss_new_client()
{
    gss_client_state *state;
    state = (gss_client_state *) malloc(sizeof(gss_client_state));
    assert(state && "Could not allocate memory for client state.");
    state->server_name = GSS_C_NO_NAME;
    state->mech_oid = GSS_C_NO_OID;
    state->context = GSS_C_NO_CONTEXT;
    state->gss_flags = GSS_C_MUTUAL_FLAG | GSS_C_SEQUENCE_FLAG;
    state->client_creds = GSS_C_NO_CREDENTIAL;
    state->username = NULL;
    state->response = NULL;
    state->response_len = 0;
    return state;
}

int auth_gss_client_init(void *s,
        const char *service,
        const char *principal, 
        const char *password, 
        const char *keytab_file, 
        uint32_t gss_flags)
{
    gss_client_state *state = (gss_client_state*) s;
    assert(state && "Client state is NULL.");
    gss_key_value_element_desc kv_element;
    gss_key_value_set_desc kv_store;
    gss_buffer_desc name_token = GSS_C_EMPTY_BUFFER;
    gss_OID mech;
    gss_name_t name = GSS_C_NO_NAME;

    if (gss_flags != 0) {
        state->gss_flags = gss_flags;
    }

    if (!(service && *service)) {
        state->maj_stat = -1;
        state->min_stat = 0;
        state->err_msg = "Must define service.";
        return AUTH_GSS_ERROR;
    }
    // Import server name first
    name_token.length = strlen(service);
    name_token.value = (char *)service;

    // could be in principal name format, i.e. service/fqdn@REALM
    if (strchr(service, '/'))
        mech = GSS_C_NO_OID;
    else
        mech = gss_krb5_nt_service_name;

    state->maj_stat = gss_import_name(&state->min_stat, 
        &name_token, mech, &state->server_name);
    if (GSS_ERROR(state->maj_stat)) {
        state->err_msg = "Could not import service name.";
        return AUTH_GSS_ERROR;
    }

    // Get credential for principal
    if (principal && *principal)
    {
        gss_buffer_desc principal_token = GSS_C_EMPTY_BUFFER;
        principal_token.length = strlen(principal);
        principal_token.value = (char *)principal;

        state->maj_stat = gss_import_name(&state->min_stat,
            &principal_token, GSS_C_NT_USER_NAME, &name);
        if (GSS_ERROR(state->maj_stat)) {
            state->err_msg = "Could not import principal.";
            return AUTH_GSS_ERROR;
        }
    }

    if(keytab_file && *keytab_file) {
        kv_element = (gss_key_value_element_desc) { 
            .key = "client_keytab", 
            .value = keytab_file
        };
    } else if (password && *password) {
        kv_element = (gss_key_value_element_desc) { 
            .key = "password", 
            .value = password 
        };
    } else {
        kv_element = (gss_key_value_element_desc) { 
            .key = "ccache", 
            .value = NULL
        };
    }
    kv_store = (gss_key_value_set_desc) { 
        .count = 1,
        .elements = &kv_element
    };

    state->maj_stat = gss_acquire_cred_from(&state->min_stat,
                                name,
                                GSS_C_INDEFINITE,
                                GSS_C_NO_OID_SET,
                                GSS_C_INITIATE,
                                &kv_store,
                                &state->client_creds,
                                NULL,
                                NULL);
    if (GSS_ERROR(state->maj_stat)) {
        state->err_msg = "Invalid Credentials.";
        return AUTH_GSS_ERROR;
    }

    state->maj_stat = gss_release_name(&state->min_stat, &name);
    if (GSS_ERROR(state->maj_stat)) {
        state->err_msg = "Could not release service name.";
        return AUTH_GSS_ERROR;
    }

    return AUTH_GSS_COMPLETE;
}

int auth_ggs_client_clean(void *s)
{
    gss_client_state *state = (gss_client_state *)s;
    OM_uint32 maj_stat;
    OM_uint32 min_stat;
    int ret = AUTH_GSS_COMPLETE;
    assert(state != NULL);

    if (state->context != GSS_C_NO_CONTEXT)
        state->maj_stat = gss_delete_sec_context(&state->min_stat, &state->context, GSS_C_NO_BUFFER);
    if (state->server_name != GSS_C_NO_NAME)
        state->maj_stat = gss_release_name(&state->min_stat, &state->server_name);
    if (state->client_creds != GSS_C_NO_CREDENTIAL)
        state->maj_stat = gss_release_cred(&state->min_stat, &state->client_creds);
    if (state->username != NULL)
    {
        free(state->username);
        state->username = NULL;
    }
    if (state->response != NULL)
    {
        free(state->response);
        state->response_len = 0;
        state->response = NULL;
    }
    if (state->maj_stat != 0) {
        return -1;
    }

    return ret;
}

int auth_gss_client_step(void *s, char* challenge, size_t challenge_len, struct gss_channel_bindings_struct* channel_bindings)
{
    gss_client_state* state = (gss_client_state*) s;
    OM_uint32 ret_flags; // Not used, but may be necessary for gss call.
    gss_buffer_desc input_token = GSS_C_EMPTY_BUFFER;
    gss_buffer_desc output_token = GSS_C_EMPTY_BUFFER;
    int ret = AUTH_GSS_CONTINUE;
    assert(state != NULL);

    // Always clear out the old response
    if (state->response != NULL)
    {
        free(state->response);
        state->response_len = 0;
        state->response = NULL;
    }

    // If there is a challenge (data from the server) we need to give it to GSS
    if (challenge && *challenge)
    {
        input_token.value = challenge;
        input_token.length = challenge_len;
    }

    // Do GSSAPI step
    state->maj_stat = gss_init_sec_context(&state->min_stat,
                                    state->client_creds,
                                    &state->context,
                                    state->server_name,
                                    state->mech_oid,
                                    (OM_uint32)state->gss_flags,
                                    0,
                                    channel_bindings,
                                    &input_token,
                                    NULL,
                                    &output_token,
                                    &ret_flags,
                                    NULL);

    if ((state->maj_stat != GSS_S_COMPLETE) && (state->maj_stat != GSS_S_CONTINUE_NEEDED))
    {
        state->err_msg = "Failed Init Sec Context";
        ret = AUTH_GSS_ERROR;
        goto end;
    }

    ret = (state->maj_stat == GSS_S_COMPLETE) ? AUTH_GSS_COMPLETE : AUTH_GSS_CONTINUE;
    // Grab the client response to send back to the server
    if (output_token.length)
    {
        state->response_len = output_token.length;
        state->response = malloc(state->response_len);
        memcpy(state->response, output_token.value, state->response_len);
        state->maj_stat = gss_release_buffer(&state->min_stat, &output_token);
    }

    // Try to get the user name if we have completed all GSS operations
    if (ret == AUTH_GSS_COMPLETE)
    {
        gss_buffer_desc name_token;
        gss_name_t gssuser = GSS_C_NO_NAME;
        state->maj_stat = gss_inquire_context(&state->min_stat, state->context, &gssuser, NULL, NULL, NULL,  NULL, NULL, NULL);
        if (GSS_ERROR(state->maj_stat))
        {
            ret = AUTH_GSS_ERROR;
            goto end;
        }

        name_token.length = 0;
        state->maj_stat = gss_display_name(&state->min_stat, gssuser, &name_token, NULL);
        if (GSS_ERROR(state->maj_stat))
        {
            if (name_token.value)
                gss_release_buffer(&state->min_stat, &name_token);
            gss_release_name(&state->min_stat, &gssuser);

            ret = AUTH_GSS_ERROR;
            goto end;
        }
        else
        {
            state->username = (char *)malloc(name_token.length + 1);
            strncpy(state->username, (char*) name_token.value, name_token.length);
            state->username[name_token.length] = 0;
            gss_release_buffer(&state->min_stat, &name_token);
            gss_release_name(&state->min_stat, &gssuser);
        }
    }
end:
    if (output_token.value)
        gss_release_buffer(&state->min_stat, &output_token);
    return ret;
}

char * auth_gss_client_response(void *s, size_t *response_len)
{
    gss_client_state* state = (gss_client_state*) s;
    assert(state != NULL);
    *response_len = state->response_len;
    return state->response;
}

int auth_gss_client_wrap_iov(void *s, char *data, size_t data_len, int protect, int *pad_len)
{
    gss_client_state* state = (gss_client_state*) s;
    assert(state != NULL);
    int iov_count = 3;
    gss_iov_buffer_desc iov[iov_count];
    int ret = AUTH_GSS_CONTINUE;
    int conf_state;

    // Always clear out the old response
    if (state->response != NULL)
    {
        free(state->response);
        state->response_len = 0;
        state->response = NULL;
    }

    iov[0].type = GSS_IOV_BUFFER_TYPE_HEADER | GSS_IOV_BUFFER_FLAG_ALLOCATE;
    iov[1].type = GSS_IOV_BUFFER_TYPE_DATA;
    iov[1].buffer.value = data;
    iov[1].buffer.length = data_len;
    iov[2].type = GSS_IOV_BUFFER_TYPE_PADDING | GSS_IOV_BUFFER_FLAG_ALLOCATE;

    state->maj_stat = gss_wrap_iov(&state->min_stat,        /* minor_status */
                         state->context,         /* context_handle */
                         protect,       /* conf_req_flag */
                         GSS_C_QOP_DEFAULT, /* qop_req */
                         &conf_state,          /* conf_state */
                         iov,           /* iov */
                         iov_count);    /* iov_count */
    if (state->maj_stat != GSS_S_COMPLETE)
    {
        ret = AUTH_GSS_ERROR;
    }
    else
    {
        ret = AUTH_GSS_COMPLETE;

        int index = 0;
        OM_uint32 stoken_len= 0;
        int bufsize = iov[0].buffer.length+
                      iov[1].buffer.length+
                      iov[2].buffer.length+
                      sizeof(unsigned int);
        char * response = (char*)malloc(bufsize);
        memset(response,0,bufsize);
        /******************************************************
        Per Microsoft 2.2.9.1.2.2.2 for kerberos encrypted data
        First section of data is a 32-bit unsigned int containing
        the length of the Security Token followed by the encrypted message.
        Encrypted data = |32-bit unsigned int|Message|
        The message must start with the security token, followed by
        the actual encrypted message.
        Message = |Security Token|encrypted data|padding
        iov[0] = security token
        iov[1] = encrypted message
        iov[2] = padding
        ******************************************************/
        /* Security Token length */
        stoken_len = iov[0].buffer.length;
        memcpy(response, &stoken_len, sizeof(unsigned int));
        index += sizeof(unsigned int);
        /* Security Token */
        memcpy(response+index, iov[0].buffer.value, iov[0].buffer.length);
        index += iov[0].buffer.length;
        /* Message */
        memcpy(response+index, iov[1].buffer.value, iov[1].buffer.length);
        index += iov[1].buffer.length;
        /* Padding */
        *pad_len = iov[2].buffer.length;
        if (*pad_len > 0)
        {
            memcpy(response+index, iov[2].buffer.value, iov[2].buffer.length);
            index += iov[2].buffer.length;
        }

        state->responseConf = conf_state;
        state->response = response;
        state->response_len = index;
    }
    (void)gss_release_iov_buffer(&state->min_stat, iov, iov_count);
    return ret;
}

int auth_gss_client_unwrap_iov(void *s, char *data, size_t len)
{
    gss_client_state* state = (gss_client_state*) s;
    assert(state != NULL);
    int conf_state = 1;
    OM_uint32 qop_state = 0;
    int ret = AUTH_GSS_COMPLETE;
    int iov_count = 3;
    gss_iov_buffer_desc iov[iov_count];
    unsigned int token_len = 0;

    // Always clear out the old response
    if (state->response != NULL)
    {
        free(state->response);
        state->response = NULL;
        state->responseConf = 0;
    }

    if (!data || len == 0)
    {
        // nothing to do, return
        data = (unsigned char *)malloc(1);
        data[0] = 0;
        state->response = (char*)data;
        return AUTH_GSS_COMPLETE;
    }

    memcpy(&token_len, data, sizeof(unsigned int));

    if (len-4-token_len < 0)
    {
        state->err_msg = "Data length error in response";
        return AUTH_GSS_ERROR;
    }

    iov[0].type = GSS_IOV_BUFFER_TYPE_HEADER;
    iov[0].buffer.value = data+4;
    iov[0].buffer.length = token_len;

    iov[1].type = GSS_IOV_BUFFER_TYPE_DATA;
    iov[1].buffer.value = data+4+token_len;
    iov[1].buffer.length = len-4-token_len;

    iov[2].type = GSS_IOV_BUFFER_TYPE_DATA;
    iov[2].buffer.value = "";
    iov[2].buffer.length = 0;

    state->maj_stat = gss_unwrap_iov(&state->min_stat, state->context, &conf_state, &qop_state, iov, iov_count);

    if (state->maj_stat != GSS_S_COMPLETE)
    {
        state->err_msg = "Unable to decode message";
        ret = AUTH_GSS_ERROR;
    }
    else
    {
        ret = AUTH_GSS_COMPLETE;

        // Grab the client response
        state->responseConf = conf_state;
        state->response_len = iov[1].buffer.length;
        state->response = malloc(state->response_len);
        assert(state->response);
        memcpy(state->response, iov[1].buffer.value, state->response_len);
    }
    return ret;
}

void *channel_bindings(OM_uint32 initiator_addrtype,
        size_t initiator_address_length,
        void *initiator_address_value,
        OM_uint32 acceptor_addrtype,
        size_t  acceptor_address_length,
        void *acceptor_address_value,
        size_t application_data_length,
        void *application_data_value)
{
    gss_channel_bindings_t cb;
    void *temp;

    cb = malloc(sizeof(gss_channel_bindings_t*));
    assert(cb && "Unable to allocate memory for channel binding.");


    if(initiator_address_length) {
        temp = malloc(initiator_address_length);
        assert(temp && "Could not reserve memory for initiator address.");
    } else {
        temp = NULL;
    }
    cb->initiator_addrtype = initiator_addrtype;
    cb->initiator_address = (gss_buffer_desc) {
        .length = initiator_address_length,
        .value = temp
    };

    if(acceptor_address_length) {
        temp = malloc(acceptor_address_length);
        assert(temp && "Could not reserve memory for acceptor address.");
    } else {
        temp = NULL;
    }
    cb->acceptor_addrtype = acceptor_addrtype;
    cb->acceptor_address = (gss_buffer_desc) {
        .length = acceptor_address_length,
        .value = temp
    };

    if(application_data_length) {
        temp = malloc(application_data_length);
        assert(temp && "Could not reserve memory for application data.");
    } else {
        temp = NULL;
    }
    cb->application_data = (gss_buffer_desc) {
        .length = application_data_length,
        .value = temp
    };
    return cb;
}

void * auth_gss_new_server()
{
    gss_server_state *state;
    state = (gss_server_state*) malloc(sizeof(gss_server_state));
    assert(state && "Could not allocate memory for server state.");
    state->context = GSS_C_NO_CONTEXT;
    state->server_name = GSS_C_NO_NAME;
    state->client_name = GSS_C_NO_NAME;
    state->server_creds = GSS_C_NO_CREDENTIAL;
    state->client_creds = GSS_C_NO_CREDENTIAL;
    state->username = NULL;
    state->targetname = NULL;
    state->response = NULL;
    state->response_len = 0;
    state->maj_stat = 0;
    state->min_stat = 0;
    state->err_msg = "";
    return state;
}

int auth_gss_server_init(void *s, char *service,
        const char *password, 
        const char *keytab_file)
{
    gss_server_state *state = (gss_server_state*) s;
    assert(state && "Server state is NULL.");
    gss_buffer_desc name_token = GSS_C_EMPTY_BUFFER;
    gss_key_value_element_desc kv_element;
    gss_key_value_set_desc kv_store;

    // Get credential for service
    if (service && *service)
    {
        name_token.length = strlen(service);
        name_token.value = (char *)service;

        state->maj_stat = gss_import_name(&state->min_stat,
            &name_token, GSS_C_NT_HOSTBASED_SERVICE, &state->server_name);
        if (GSS_ERROR(state->maj_stat)) {
            state->err_msg = "Could not import service.";
            return AUTH_GSS_ERROR;
        }
    }

    if(keytab_file && *keytab_file) {
        kv_element = (gss_key_value_element_desc) { 
            .key = "client_keytab", 
            .value = keytab_file
        };
    } else if (password && *password) {
        kv_element = (gss_key_value_element_desc) { 
            .key = "password", 
            .value = password 
        };
    } else {
        kv_element = (gss_key_value_element_desc) { 
            .key = "ccache", 
            .value = NULL
        };
    }
    kv_store = (gss_key_value_set_desc) { 
        .count = 1,
        .elements = &kv_element
    };

    state->maj_stat = gss_acquire_cred_from(&state->min_stat,
                                state->server_name,
                                GSS_C_INDEFINITE,
                                GSS_C_NO_OID_SET,
                                GSS_C_ACCEPT,
                                &kv_store,
                                &state->server_creds,
                                NULL,
                                NULL);
    if (GSS_ERROR(state->maj_stat)) {
        state->err_msg = "Invalid Credentials.";
        return AUTH_GSS_ERROR;
    }

    return AUTH_GSS_COMPLETE;    
}

int auth_gss_server_clean(void *s)
{
    gss_server_state *state = (gss_server_state*) s;
    assert(state && "Server state is NULL.");
    int ret = AUTH_GSS_COMPLETE;

    if (state->context != GSS_C_NO_CONTEXT)
        gss_delete_sec_context(&state->min_stat, &state->context, GSS_C_NO_BUFFER);
    if (state->server_name != GSS_C_NO_NAME)
        gss_release_name(&state->min_stat, &state->server_name);
    if (state->client_name != GSS_C_NO_NAME)
        gss_release_name(&state->min_stat, &state->client_name);
    if (state->server_creds != GSS_C_NO_CREDENTIAL)
        gss_release_cred(&state->min_stat, &state->server_creds);
    if (state->client_creds != GSS_C_NO_CREDENTIAL)
        gss_release_cred(&state->min_stat, &state->client_creds);
    if (state->username != NULL)
    {
        free(state->username);
        state->username = NULL;
    }
    if (state->targetname != NULL)
    {
        free(state->targetname);
        state->targetname = NULL;
    }
    if (state->response != NULL)
    {
        free(state->response);
        state->response = NULL;
        state->response_len = 0;
    }

    return ret;    
}

int auth_gss_server_step(void *s, char *challenge, size_t challenge_len)
{
    gss_server_state *state = (gss_server_state*) s;
    assert(state && "Server state is NULL.");
    gss_buffer_desc input_token = GSS_C_EMPTY_BUFFER;
    gss_buffer_desc output_token = GSS_C_EMPTY_BUFFER;
    gss_name_t target_name = GSS_C_NO_NAME;
    int ret = AUTH_GSS_CONTINUE;

    // Always clear out the old response
    if (state->response != NULL)
    {
        free(state->response);
        state->response = NULL;
        state->response_len = 0;
    }

    // If there is a challenge (data from the server) we need to give it to GSS
    if (challenge && *challenge)
    {
        size_t len;
        input_token.value = challenge;
        input_token.length = challenge_len;
    }
    else
    {
        state->err_msg =  "No challenge parameter in request from client";
        state->maj_stat = 0;
        state->min_stat = 0;
        ret = AUTH_GSS_ERROR;
        goto end;
    }

    state->maj_stat = gss_accept_sec_context(&state->min_stat,
                                      &state->context,
                                      state->server_creds,
                                      &input_token,
                                      GSS_C_NO_CHANNEL_BINDINGS,
                                      &state->client_name,
                                      NULL,
                                      &output_token,
                                      NULL,
                                      NULL,
                                      &state->client_creds);

    if (GSS_ERROR(state->maj_stat))
    {
        state->err_msg = "Failed Accept Sec Context.";
        ret = AUTH_GSS_ERROR;
        goto end;
    }

    // Grab the server response to send back to the client
    if (output_token.length)
    {
        state->response_len = output_token.length;
        state->response = malloc(state->response_len);
        assert(state->response && "Could not allocate memory for response token.");
        memcpy(state->response, output_token.value, state->response_len);
        state->maj_stat = gss_release_buffer(&state->min_stat, &output_token);
        if(GSS_ERROR(state->maj_stat)) {
            state->err_msg = "Error releasing output token buffer.";
            goto end;
        }
    }

    // Get the user name
    state->maj_stat = gss_display_name(&state->min_stat, state->client_name, &output_token, NULL);
    if (GSS_ERROR(state->maj_stat))
    {
        state->err_msg = "Could not get user name from server state.";
        ret = AUTH_GSS_ERROR;
        goto end;
    }
    state->username = (char *)malloc(output_token.length + 1);
    assert(state->username && "Could not allocate memory for username from output token.");
    strncpy(state->username, (char*) output_token.value, output_token.length);
    state->username[output_token.length] = 0;

    // Get the target name if no server creds were supplied
    if (state->server_creds == GSS_C_NO_CREDENTIAL)
    {
        state->maj_stat = gss_inquire_context(&state->min_stat, 
            state->context, 
            NULL, 
            &target_name, 
            NULL, 
            NULL, 
            NULL, 
            NULL, 
            NULL);
        if (GSS_ERROR(state->maj_stat))
        {
            state->err_msg = "Could not extract target name from server context.";
            ret = AUTH_GSS_ERROR;
            goto end;
        }

        // Free output token if necessary before reusing
        if (output_token.length)
            state->maj_stat = gss_release_buffer(&state->min_stat, &output_token);
        if(GSS_ERROR(state->maj_stat)) {
            state->err_msg = "Could not release output token.";
            goto end;
        }

        state->maj_stat = gss_display_name(&state->min_stat, target_name, &output_token, NULL);
        if (GSS_ERROR(state->maj_stat))
        {
            state->err_msg = "Could not get target name from output token.";
            ret = AUTH_GSS_ERROR;
            goto end;
        }
        state->targetname = (char *)malloc(output_token.length + 1);
        assert(state->targetname && "Could not allocate memory for target name.");
        strncpy(state->targetname, (char*) output_token.value, output_token.length);
        state->targetname[output_token.length] = 0;
    }

    ret = AUTH_GSS_COMPLETE;

end:
    if (target_name != GSS_C_NO_NAME)
        gss_release_name(&state->min_stat, &target_name);
    if (output_token.length)
        gss_release_buffer(&state->min_stat, &output_token);
    return ret;    
}

void debug_gss_client_state(void *s)
{
    gss_client_state *state = (gss_client_state*) s;
    assert(state && "Server state is NULL.");
    fprintf(stderr, "gss=%p\n", state);
    fprintf(stderr, "   context=%p\n", state->context);
    fprintf(stderr, "   server_name=%p\n", state->server_name);
    fprintf(stderr, "   mech_oid=%p\n", state->mech_oid);
    fprintf(stderr, "   gss_flags=%ld\n", state->gss_flags);
    fprintf(stderr, "   client_creds=%p\n", state->client_creds);
    fprintf(stderr, "   username=%s\n", state->username);
    fprintf(stderr, "   response=%s\n", state->response);
    fprintf(stderr, "   response_len=%ld\n", state->response_len);
    fprintf(stderr, "   responseConf=%d\n", state->responseConf);
    fprintf(stderr, "   maj_stat=%d\n", state->maj_stat);
    fprintf(stderr, "   min_stat=%d\n", state->min_stat);
    fprintf(stderr, "   err_msg=%s\n", state->err_msg);

}