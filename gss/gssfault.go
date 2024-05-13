package gss

import "C"

import (
	"fmt"
)

var CallingErrors = map[int32]string {
	1: "GSS_S_CALL_INACCESSIBLE_READ",
	2: "GSS_S_CALL_INACCESSIBLE_WRITE",
	3: "GSS_S_CALL_BAD_STRUCTURE",
}

var RoutineErrors = map[int32]string {
	1: "GSS_S_BAD_MECH",
	2: "GSS_S_BAD_NAME",
	3: "GSS_S_BAD_NAMETYPE",
	4: "GSS_S_BAD_BINDINGS",
	5: "GSS_S_BAD_STATUS",
	6: "GSS_S_BAD_MIC, GSS_S_BAD_SIG",
	7: "GSS_S_NO_CRED",
	8: "GSS_S_NO_CONTEXT",
	9: "GSS_S_DEFECTIVE_TOKEN",
	10: "GSS_S_DEFECTIVE_CREDENTIAL",
	11: "GSS_S_CREDENTIALS_EXPIRED",
	12: "GSS_S_CONTEXT_EXPIRED",
	13: "GSS_S_FAILURE",
	14: "GSS_S_BAD_QOP",
	15: "GSS_S_UNAUTHORIZED",
	16: "GSS_S_UNAVAILABLE",
	17: "GSS_S_DUPLICATE_ELEMENT",
	18: "GSS_S_NAME_NOT_MN",
}

var SupplementaryInfo = map[int32]string {
	1: "GSS_S_CONTINUE_NEEDED",
	2: "GSS_S_DUPLICATE_TOKEN",
	4: "GSS_S_OLD_TOKEN",
	8: "GSS_S_UNSEQ_TOKEN",
	16: "GSS_S_GAP_TOKEN",
}

var Krb5StatusCodes = map[int32]string {
	-1765328384: "KRB5KDC_ERR_NONE",
	-1765328383: "KRB5KDC_ERR_NAME_EXP",
	-1765328382: "KRB5KDC_ERR_SERVICE_EXP",
	-1765328381: "KRB5KDC_ERR_BAD_PVNO",
	-1765328380: "KRB5KDC_ERR_C_OLD_MAST_KVNO",
	-1765328379: "KRB5KDC_ERR_S_OLD_MAST_KVNO",
	-1765328378: "KRB5KDC_ERR_C_PRINCIPAL_UNKNOWN",
	-1765328377: "KRB5KDC_ERR_S_PRINCIPAL_UNKNOWN",
	-1765328376: "KRB5KDC_ERR_PRINCIPAL_NOT_UNIQUE",
	-1765328375: "KRB5KDC_ERR_NULL_KEY",
	-1765328374: "KRB5KDC_ERR_CANNOT_POSTDATE",
	-1765328373: "KRB5KDC_ERR_NEVER_VALID",
	-1765328372: "KRB5KDC_ERR_POLICY",
	-1765328371: "KRB5KDC_ERR_BADOPTION",
	-1765328370: "KRB5KDC_ERR_ETYPE_NOSUPP",
	-1765328369: "KRB5KDC_ERR_SUMTYPE_NOSUPP",
	-1765328368: "KRB5KDC_ERR_PADATA_TYPE_NOSUPP",
	-1765328367: "KRB5KDC_ERR_TRTYPE_NOSUPP",
	-1765328366: "KRB5KDC_ERR_CLIENT_REVOKED",
	-1765328365: "KRB5KDC_ERR_SERVICE_REVOKED",
	-1765328364: "KRB5KDC_ERR_TGT_REVOKED",
	-1765328363: "KRB5KDC_ERR_CLIENT_NOTYET",
	-1765328362: "KRB5KDC_ERR_SERVICE_NOTYET",
	-1765328361: "KRB5KDC_ERR_KEY_EXP",
	-1765328360: "KRB5KDC_ERR_PREAUTH_FAILED",
	-1765328359: "KRB5KDC_ERR_PREAUTH_REQUIRED",
	-1765328358: "KRB5KDC_ERR_SERVER_NOMATCH",
	-1765328352: "RB5KRB_AP_ERR_TKT_EXPIRED",
	-1765328351: "KRB5KRB_AP_ERR_TKT_NYV",
	-1765328350: "KRB5KRB_AP_ERR_REPEAT",
	-1765328349: "KRB5KRB_AP_ERR_NOT_US",
	-1765328348: "KRB5KRB_AP_ERR_BADMATCH",
	-1765328347: "KRB5KRB_AP_ERR_SKEW",
	-1765328346: "KRB5KRB_AP_ERR_BADADDR",
	-1765328345: "KRB5KRB_AP_ERR_BADVERSION",
	-1765328344: "KRB5KRB_AP_ERR_MSG_TYPE",
	-1765328343: "KRB5KRB_AP_ERR_MODIFIED",
	-1765328342: "KRB5KRB_AP_ERR_BADORDER",
	-1765328341: "KRB5KRB_AP_ERR_ILL_CR_TKT",
	-1765328340: "KRB5KRB_AP_ERR_BADKEYVER",
	-1765328339: "KRB5KRB_AP_ERR_NOKEY",
	-1765328338: "KRB5KRB_AP_ERR_MUT_FAIL",
	-1765328337: "KRB5KRB_AP_ERR_BADDIRECTION",
	-1765328336: "KRB5KRB_AP_ERR_METHOD",
	-1765328335: "KRB5KRB_AP_ERR_BADSEQ",
	-1765328334: "KRB5KRB_AP_ERR_INAPP_CKSUM",
	-1765328255: "KRB5_LIBOS_BADLOCKFLAG",
	-1765328254: "KRB5_LIBOS_CANTREADPWD",
	-1765328253: "KRB5_LIBOS_BADPWDMATCH",
	-1765328252: "KRB5_LIBOS_PWDINTR",
	-1765328251: "KRB5_PARSE_ILLCHAR",
	-1765328250: "KRB5_PARSE_MALFORMED",
	-1765328249: "KRB5_CONFIG_CANTOPEN",
	-1765328248: "KRB5_CONFIG_BADFORMAT",
	-1765328247: "KRB5_CONFIG_NOTENUFSPACE",
	-1765328246: "KRB5_BADMSGTYPE",
	-1765328245: "KRB5_CC_BADNAME",
	-1765328244: "KRB5_CC_UNKNOWN_TYPE",
	-1765328243: "KRB5_CC_NOTFOUND",
	-1765328242: "KRB5_CC_END",
	-1765328241: "KRB5_NO_TKT_SUPPLIED",
	-1765328240: "KRB5KRB_AP_WRONG_PRINC",
	-1765328239: "KRB5KRB_AP_ERR_TKT_INVALID",
	-1765328238: "KRB5_PRINC_NOMATCH",
	-1765328237: "KRB5_KDCREP_MODIFIED",
	-1765328236: "KRB5_KDCREP_SKEW",
	-1765328235: "KRB5_IN_TKT_REALM_MISMATCH",
	-1765328234: "KRB5_PROG_ETYPE_NOSUPP",
	-1765328233: "KRB5_PROG_KEYTYPE_NOSUPP",
	-1765328232: "KRB5_WRONG_ETYPE",
	-1765328231: "KRB5_PROG_SUMTYPE_NOSUPP",
	-1765328230: "KRB5_REALM_UNKNOWN",
	-1765328229: "KRB5_SERVICE_UNKNOWN",
	-1765328228: "KRB5_KDC_UNREACH",
	-1765328227: "KRB5_NO_LOCALNAME",
	-1765328226: "KRB5_MUTUAL_FAILED",
	-1765328225: "KRB5_RC_TYPE_EXISTS",
	-1765328224: "KRB5_RC_MALLOC",
	-1765328223: "KRB5_RC_TYPE_NOTFOUND",
	-1765328222: "KRB5_RC_UNKNOWN",
	-1765328221: "KRB5_RC_REPLAY",
	-1765328220: "KRB5_RC_IO",
	-1765328219: "KRB5_RC_NOIO",
	-1765328218: "KRB5_RC_PARSE",
	-1765328217: "KRB5_RC_IO_EOF",
	-1765328216: "KRB5_RC_IO_MALLOC",
	-1765328215: "KRB5_RC_IO_PERM",
	-1765328214: "KRB5_RC_IO_IO",
	-1765328213: "KRB5_RC_IO_UNKNOWN",
	-1765328212: "KRB5_RC_IO_SPACE",
	-1765328211: "KRB5_TRANS_CANTOPEN",
	-1765328210: "KRB5_TRANS_BADFORMAT",
	-1765328209: "KRB5_LNAME_CANTOPEN",
	-1765328208: "KRB5_LNAME_NOTRANS",
	-1765328207: "KRB5_LNAME_BADFORMAT",
	-1765328206: "KRB5_CRYPTO_INTERNAL",
	-1765328205: "KRB5_KT_BADNAME",
	-1765328204: "KRB5_KT_UNKNOWN_TYPE",
	-1765328203: "KRB5_KT_NOTFOUND",
	-1765328202: "KRB5_KT_END",
	-1765328201: "KRB5_KT_NOWRITE",
	-1765328200: "KRB5_KT_IOERR",
	-1765328199: "KRB5_NO_TKT_IN_RLM",
	-1765328198: "KRB5DES_BAD_KEYPAR",
	-1765328197: "KRB5DES_WEAK_KEY",
	-1765328196: "KRB5_BAD_ENCTYPE",
	-1765328195: "KRB5_BAD_KEYSIZE",
	-1765328194: "KRB5_BAD_MSIZE",
	-1765328193: "KRB5_CC_TYPE_EXISTS",
	-1765328192: "KRB5_KT_TYPE_EXISTS",
	-1765328191: "KRB5_CC_IO",
	-1765328190: "KRB5_FCC_PERM",
	-1765328189: "KRB5_FCC_NOFILE",
	-1765328188: "KRB5_FCC_INTERNAL",
	-1765328187: "KRB5_CC_WRITE",
	-1765328186: "KRB5_CC_NOMEM",
	-1765328185: "KRB5_CC_FORMAT",
	-1765328184: "KRB5_INVALID_FLAGS",
	-1765328183: "KRB5_NO_2ND_TKT",
	-1765328182: "KRB5_NOCREDS_SUPPLIED",
	-1765328181: "KRB5_SENDAUTH_BADAUTHVERS",
	-1765328180: "KRB5_SENDAUTH_BADAPPLVERS",
	-1765328179: "KRB5_SENDAUTH_BADRESPONSE",
	-1765328178: "KRB5_SENDAUTH_REJECTED",
	-1765328177: "KRB5_PREAUTH_BAD_TYPE",
	-1765328176: "KRB5_PREAUTH_NO_KEY",
	-1765328175: "KRB5_PREAUTH_FAILED",
	-1765328174: "KRB5_RCACHE_BADVNO",
	-1765328173: "KRB5_CCACHE_BADVNO",
	-1765328172: "KRB5_KEYTAB_BADVNO",
	-1765328171: "KRB5_PROG_ATYPE_NOSUPP",
	-1765328170: "KRB5_RC_REQUIRED",
	-1765328169: "KRB5_ERR_BAD_HOSTNAME",
	-1765328168: "KRB5_ERR_HOST_REALM_UNKNOWN",
	-1765328167: "KRB5_SNAME_UNSUPP_NAMETYPE",
	-1765328166: "KRB5KRB_AP_ERR_V4_REPLY",
	-1765328165: "KRB5_REALM_CANT_RESOLVE",
	-1765328164: "KRB5_TKT_NOT_FORWARDABLE",
	-1765328163: "KRB5_FWD_BAD_PRINCIPAL",
	-1765328162: "KRB5_GET_IN_TKT_LOOP",
	-1765328161: "KRB5_CONFIG_NODEFREALM",
	-1765328160: "KRB5_SAM_UNSUPPORTED",
	-1765328159: "KRB5_KT_NAME_TOOLONG",
	-1765328158: "KRB5_KT_KVNONOTFOUND",
	-1765328157: "KRB5_CONF_NOT_CONFIGURED",
}

type GssFault struct {
    MajStat C.uint
    MinStat C.uint
    Message string
    Func string
    Status C.int
}

func (gf GssFault) Error () string {
    var supplementary_info int32
    var routine_error int32
    var calling_error int32
    supplementary_info = int32(gf.MajStat & 0xffff)
    routine_error = int32((gf.MajStat >> 16) & 0xff)
    calling_error = int32((gf.MajStat >> 24) & 0xff)
    return fmt.Sprintf("%s: supplementary_info=%d(%s) routine_error=%d(%s) calling_error=%d(%s) min=%d(%s) %s",
        gf.Func,
        supplementary_info, SupplementaryInfo[supplementary_info], 
        routine_error, RoutineErrors[routine_error],
        calling_error, CallingErrors[calling_error],
        int32(gf.MinStat), Krb5StatusCodes[int32(gf.MinStat)],
        gf.Message)
}
