package transport

import (
	"net/http"
    "github.com/samanamonitor/gosammwr/sammwr/gss"
    "net/url"
    "encoding/base64"
    "fmt"
    "strings"
    "bytes"
    "io"
    "errors"
)

func multipart_encode(f *bytes.Buffer, encrypted_message []byte, message_len int) {
    f.Write([]byte(fmt.Sprint("--Encrypted Boundary\r\n", 
        "Content-Type: application/HTTP-SPNEGO-session-encrypted\r\n",
        "OriginalContent: type=application/soap+xml;charset=UTF-8;Length=", message_len, "\r\n",
        "--Encrypted Boundary\r\n",
        "Content-Type: application/octet-stream\r\n")))
    f.Write(encrypted_message)
    f.Write([]byte("--Encrypted Boundary--\r\n"))
}

func multipart_decode(data []byte) []byte {
    elements := bytes.Split(data, []byte("\r\n"))
    separator := elements[0][:len(elements[0])-2]
    encrypted_data := bytes.Split(elements[5], separator)
    return encrypted_data[0]
}

type Transport struct {
	endpoint_string string
	endpoint *url.URL
	username string
	password string
	keytab_file string
	authenticated bool
	client *http.Client
	gssAuth gss.Gss
	service string
	challenge []byte

}

func (self *Transport) Init(
		endpoint string,
		username string,
		password string,
		keytab_file string) error {
	var err error
	self.endpoint_string = endpoint
	self.endpoint, err = url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	self.username = username
	self.password = password
	self.keytab_file = keytab_file

	self.client = &http.Client{}
	self.service = "HTTP"
	self.gssAuth = gss.Gss{}
	result := self.gssAuth.AuthGssClientInit(self.service + "/" + self.endpoint.Hostname(), 
		self.username, self.password, self.keytab_file, 0)
	if result != gss.AUTH_GSS_COMPLETE {
		return errors.New(fmt.Sprintf("Unable to initialize GSS as client. Maj=%08x Min=%d",
			self.gssAuth.Maj_stat, self.gssAuth.Min_stat))
	}
	return nil
}

func (self *Transport) prepareRequest(message io.Reader) (*http.Request, error) {
    req, err := http.NewRequest("POST", self.endpoint_string, message)
    req.Header.Add("Accept-Encoding", "gzip, deflate")
    req.Header.Add("Accept", "*.*")
    req.Header.Add("Connection", "Keep-Alive")
    req.Header.Add("Content-Type", "multipart/encrypted;protocol=\"application/HTTP-SPNEGO-session-encrypted\";boundary=\"Encrypted Boundary\"");

    return req, err
}

func (self *Transport) BuildSession() error {
	var result int
	var challenge []byte

	result = gss.AUTH_GSS_CONTINUE
	for ; result == gss.AUTH_GSS_CONTINUE; {

		result = self.gssAuth.AuthGssClientStep(challenge)
        if result == gss.AUTH_GSS_ERROR {
        	return errors.New(fmt.Sprintf("AuthGssClientStep failed with maj=%08x min=%d",
        		self.gssAuth.Maj_stat, int32(self.gssAuth.Min_stat)))
        }

        if result == gss.AUTH_GSS_CONTINUE {
            challenge = self.gssAuth.AuthGssClientResponse()

            /* send to server */
            challenge_b64 := base64.StdEncoding.EncodeToString(challenge)

            req, err := self.prepareRequest(bytes.NewBuffer(nil))
            if err != nil {
            	return err
            }
	        req.Header.Add("Authorization", "Negotiate " + challenge_b64)
            resp, err := self.client.Do(req)
            if err != nil {
                return err
            }
            if resp.StatusCode != 200 {
            	return errors.New(fmt.Sprintf("Error. StatusCode=%d", resp.StatusCode))
            }
            authentication_header := resp.Header["Www-Authenticate"]
            if len(authentication_header) < 1 {
                return errors.New(fmt.Sprintf("Error. Invalid Www-Authenticate header. %s", authentication_header))
            }
            temp := strings.Split(authentication_header[0], " ")
            if len(temp) < 2 {
                return errors.New(fmt.Sprintf("Error. Invalid authentication token. %s", authentication_header[0]))
            }
            challenge_b64 = temp[1]
            challenge, _ = base64.StdEncoding.DecodeString(challenge_b64)
        }		
	}
	self.authenticated = true
	return nil
}

func (self *Transport) SendMessage(message []byte) (error, []byte) {
	var err error

	if self.authenticated == false {
		err = self.BuildSession()
		if err != nil {
			return err, []byte{}
		}
	}
    result := self.gssAuth.AuthGSSClientWrapIov(message)
    if result != gss.AUTH_GSS_COMPLETE {
    	return errors.New(fmt.Sprintf("AuthGSSClientWrapIov failed with maj=%08x min=%d",
    		self.gssAuth.Maj_stat, self.gssAuth.Min_stat)), []byte{}
    }
    f := bytes.NewBuffer(nil)
    multipart_encode(f, self.gssAuth.AuthGssClientResponse(), len(message))
    req, _ := self.prepareRequest(f)
    resp, _ := self.client.Do(req)
    defer resp.Body.Close()

    response_data, _ := io.ReadAll(resp.Body)
    encrypted_message := multipart_decode(response_data)
    result = self.gssAuth.AuthGSSClientUnwrapIov(encrypted_message)
    if result != gss.AUTH_GSS_COMPLETE {
    	return errors.New(fmt.Sprintf("AuthGSSClientUnwrapIov failed with maj=%08x min=%d",
    		self.gssAuth.Maj_stat, self.gssAuth.Min_stat)), []byte{}
    }
    response_clear := self.gssAuth.AuthGssClientResponse()
    if resp.StatusCode != 200 {
        return errors.New(fmt.Sprintf("Error. StatusCode=%d\n%s\n",
        	resp.StatusCode, response_clear)), []byte{}
    }

    return nil, response_clear
}

func (self *Transport) Close() error {
	result := self.gssAuth.AuthGssClientClean()
    if result != gss.AUTH_GSS_COMPLETE {
    	return errors.New(fmt.Sprintf("AuthGssClientClean failed with maj=%08x min=%d",
    		self.gssAuth.Maj_stat, self.gssAuth.Min_stat))
    }
    return nil
}
