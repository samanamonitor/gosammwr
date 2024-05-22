package transport

/*
TODO:

Move mime boundary from burnt into the code to extracting it from HTTP headers
*/

import (
	"net/http"
	"github.com/samanamonitor/gosammwr/gss"
	"net/url"
	"encoding/base64"
	"fmt"
	"strings"
	"bytes"
	"io"
	"errors"
)

type TransportFault struct {
	Err error
	StatusCode int
	Message string
	Payload []byte
}

func (tf *TransportFault) Error() string {
	return fmt.Sprintf("Transport Error: StatusCode=%d Message=%s\n%s",
		tf.StatusCode, tf.Message, tf.Err )
}

type Transport struct {
	Endpoint string
	endpoint_url *url.URL
	Username string
	Password string
	Keytab_file string
	authenticated bool
	client *http.Client
	gssAuth gss.Gss
	service string
	challenge []byte
	MultipartBoundary string
	ProtoString string
	ContentType string
	Encrypt bool
}

func NewTransport(endpoint string, username string, password string, keytab_file string) (*Transport, error) {
	var err error

	self := &Transport{
		Endpoint: endpoint,
		Username: username,
		Password: password,
		Keytab_file: keytab_file,
	}

	self.endpoint_url, err = url.Parse(self.Endpoint)
	if err != nil {
		return nil, err
	}

	self.client = &http.Client{}
	self.service = "HTTP"
	self.gssAuth = gss.Gss{}
	result := self.gssAuth.AuthGssClientInit(self.service + "/" + self.endpoint_url.Hostname(),
		self.Username, self.Password, self.Keytab_file, 0)
	if result.Status != gss.AUTH_GSS_COMPLETE {
		return nil, result
	}
	self.MultipartBoundary = "SAMM Encrypted Boundary"
	self.ProtoString = "application/HTTP-SPNEGO-session-encrypted"
	self.ContentType = "application/soap+xml;charset=UTF-8"
	self.Encrypt = true
	return self, nil
}

func (self *Transport) prepareRequest(message []byte) (*http.Request, error) {
	body := bytes.NewBuffer(nil)
	var content_type string

	if len(message) > 0 {
		if self.Encrypt {
			result := self.gssAuth.AuthGSSClientWrapIov(message)
			if result.Status != gss.AUTH_GSS_COMPLETE {
				return nil, result
			}

			encrypted_request := self.gssAuth.AuthGssClientResponse()
			mp := NewMultiPart(self.MultipartBoundary, encrypted_request)
			mp.AddHeader("Content-Type", self.ProtoString)
			mp.AddHeader("OriginalContent", fmt.Sprintf("type=%s;Length=%d", self.ContentType, len(message)))

			body = mp.Body()
			content_type = fmt.Sprint("multipart/encrypted;protocol=\"", self.ProtoString, "\"",
				";boundary=\"", self.MultipartBoundary, "\"")

		} else {
			body.Write(message)
			content_type = self.ContentType
		}
	}

	req, err := http.NewRequest("POST", self.Endpoint, body)
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept", "*.*")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Content-Type", content_type);

	return req, err
}

func extractBoundary(content_types []string) (string, error) {
	var boundary string
	for j := range content_types {
		content_type := content_types[j]
		params := strings.Split(content_type, ";")
		for i := range params {
			if strings.HasPrefix(params[i], "boundary") {
				temp := strings.Split(params[i], "=")
				if len(temp) != 2 {
					return "", errors.New("Invalid Content-Type header. Boundary missing.")
				}
				boundary = temp[1]
				if boundary[0] == '"' {
					boundary = boundary[1:]
				}
				if boundary[len(boundary)-1] == '"' {
					boundary = boundary[:len(boundary)-1]
				}
				return boundary, nil
			}
		}
	}
	return "", errors.New("Missing boundary string.")
}

func (self *Transport) processResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode == 400 {
		err := TransportFault {
			StatusCode: resp.StatusCode,
			Message: "Server couldn't decode message",
		}
		return []byte(""), &err
	}

	response_data, _ := io.ReadAll(resp.Body)
	var response_clear []byte

	if self.Encrypt {
		boundary, err := extractBoundary(resp.Header["Content-Type"])
		if err != nil {
			return []byte{}, err
		}

		_, encrypted_message, err := MultipartDecode(response_data, boundary)
		if err != nil {
			return []byte{}, err
		}
		result := self.gssAuth.AuthGSSClientUnwrapIov(encrypted_message)
		if result.Status != gss.AUTH_GSS_COMPLETE {
			return []byte{}, result
		}
		response_clear = self.gssAuth.AuthGssClientResponse()
	} else {
		response_clear = response_data
	}
	if resp.StatusCode != 200 {
		err := TransportFault {
			StatusCode: resp.StatusCode,
			Message: "Details in Payload",
			Payload: response_clear,
		}
		return response_clear, &err
	}
	return response_clear, nil
}

func (self *Transport) BuildSession() error {
	var challenge []byte

	result := gss.GssFault{ Status: gss.AUTH_GSS_CONTINUE }
	for ; result.Status == gss.AUTH_GSS_CONTINUE; {

		result = self.gssAuth.AuthGssClientStep(challenge)
		if result.Status == gss.AUTH_GSS_ERROR {
			return result
		}

		if result.Status == gss.AUTH_GSS_CONTINUE {
			challenge = self.gssAuth.AuthGssClientResponse()

			/* send to server */
			challenge_b64 := base64.StdEncoding.EncodeToString(challenge)

			req, err := self.prepareRequest([]byte(""))
			if err != nil {
				return err
			}
			req.Header.Add("Authorization", "Negotiate " + challenge_b64)
			resp, err := self.client.Do(req)
			if err != nil {
				return err
			}
			if resp.StatusCode != 200 {
				err := TransportFault{
					StatusCode: resp.StatusCode,
				}
				return &err
		   }
			authentication_header := resp.Header["Www-Authenticate"]
			if len(authentication_header) < 1 {
				err := TransportFault {
					StatusCode: resp.StatusCode,
					Message: fmt.Sprintf("Error. Invalid Www-Authenticate header. %s", authentication_header),
				}
				return &err
			}
			temp := strings.Split(authentication_header[0], " ")
			if len(temp) < 2 {
				err := TransportFault {
					StatusCode: resp.StatusCode,
					Message: fmt.Sprintf(fmt.Sprintf("Error. Invalid authentication token. %s", authentication_header[0])),
				}
				return &err
			}
			challenge_b64 = temp[1]
			challenge, _ = base64.StdEncoding.DecodeString(challenge_b64)
		}
	}

	self.authenticated = true
	return nil
}

func (self *Transport) SendMessage(message []byte) ([]byte, error) {
	var err error

	if self.authenticated == false {
		err = self.BuildSession()
		if err != nil {
			return []byte{}, err
		}
	}
	req, err := self.prepareRequest(message)
	if err != nil {
		return []byte{}, err
	}
	resp, _ := self.client.Do(req)
	defer resp.Body.Close()

	return self.processResponse(resp)
}

func (self *Transport) Close() error {
	result := self.gssAuth.AuthGssClientClean()
	if result.Status != gss.AUTH_GSS_COMPLETE {
		return result
	}
	return nil
}
