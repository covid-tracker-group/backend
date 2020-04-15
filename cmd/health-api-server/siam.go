package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Principal struct {
	Uid      string
	SelLevel int
}

func getSAMLAssertionFromHTTPHeader(r *http.Request) []byte {
	i := 1
	var assertion []byte
	for {
		value := r.Header.Get(fmt.Sprintf("X-saml-attribute-token%d", i))
		if value == "" {
			break
		}
		assertion = append(assertion, []byte(value)...)
		i++
	}
	if len(assertion) == 0 {
		return nil
	}
	decoded, err := base64.StdEncoding.DecodeString(string(assertion))
	if err != nil {
		return nil
	}
	return decoded
}

func getSamlStringAttribute(attr *SAMLAttribute) string {
	for _, value := range attr.Values {
		if value.Type == "xs:string" {
			return strings.TrimSpace(value.Value)
		}
	}
	return ""
}

func decodeSAML(r *http.Request) (*Principal, error) {
	headerValue := getSAMLAssertionFromHTTPHeader(r)
	if headerValue == nil {
		return nil, nil
	}

	samlAssertion := SAMLAssertion{}
	if err := xml.Unmarshal(headerValue, &samlAssertion); err != nil {
		return nil, err
	}

	principal := &Principal{}
	for _, attribute := range samlAssertion.Attributes {
		switch attribute.Name {
		case "uid":
			principal.Uid = getSamlStringAttribute(&attribute)
		case "sel_level":
			str := getSamlStringAttribute(&attribute)
			principal.SelLevel, _ = strconv.Atoi(str)
		}
	}
	return principal, nil
}
