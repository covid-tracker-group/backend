package main

import (
	"net/http/httptest"
	"testing"
)

func TestGetSAMLAssertionFromHTTPHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	if getSAMLAssertionFromHTTPHeader(r) != nil {
		t.Error("Missing header does not result in nil result")
	}

	r.Header.Add("X-saml-attribute-token1", "part1")
	if string(getSAMLAssertionFromHTTPHeader(r)) != "part1" {
		t.Fatal("Failed to find single-header SAML Assertion from HTTP headers")
	}

	r.Header.Add("X-saml-attribute-token2", "part2")
	r.Header.Add("X-saml-attribute-token3", "part3")
	if string(getSAMLAssertionFromHTTPHeader(r)) != "part1part2part3" {
		t.Fatal("Data from consecutive headers not read")
	}
}

func TestDecodeSAML(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	assertion, err := decodeSAML(r)
	if assertion != nil {
		t.Errorf("No assertion provided, but non-nil returned: %v", assertion)
	}
	if err != nil {
		t.Errorf("No assertion provided, but non-nil error: %v", err)
	}

	// This is the "URL decoded en base64 decoded SAML token ingeval van authenticatie via UZI- pas"
	// example from "SIAM Howto Proxy Mode v1.9"
	r.Header.Add("X-saml-attribute-token1", `
	<?xml version="1.0" encoding="UTF-8"?>
	<saml:Assertion ID="I8147B996080D693A3DFE302EAFB847D07B758D55" IssueInstant="2013-05-18T06:46:40.556Z" Version="2.0" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
		<saml:Issuer Format="urn:oasis:names:tc:SAML:2.0:nameid-format:entity" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">https://siam.anoigo.nl/aselectserver/server</saml:Issuer>
		<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
			<ds:SignedInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
				<ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#" xmlns:ds="http://www.w3.org/2000/09/ xmldsig#"/>
				<ds:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1" xmlns:ds="http://www.w3.org/2000/09/xmldsig#"/>
				<ds:Reference URI="#I8147B996080D693A3DFE302EAFB847D07B758D55" xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
					<ds:Transforms xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
						<ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature" xmlns:ds="http://www.w3.org/2000/09/ xmldsig#"/>
						<ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#" xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
							<ec:InclusiveNamespaces PrefixList="ds saml xs xsi" xmlns:ec="http://www.w3.org/2001/10/xml-exc-c14n#"/>
						</ds:Transform>
					</ds:Transforms>
					<ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1" xmlns:ds="http://www.w3.org/2000/09/xmldsig#"/>
					<ds:DigestValue xmlns:ds="http://www.w3.org/2000/09/xmldsig#">0agJc/CV2iqGJdQxC017Pfaqtyc=</ds:DigestValue>
				</ds:Reference>
			</ds:SignedInfo>
			<ds:SignatureValue xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
				NqCZXgh2/zhW7fgie7myygTB7Py2XyI1/Tnlg9VeSFcOb7wSG8IPZPKEOfvzTFcfwsC7ZKJBMKkF
				++suQKqFj7PURUKuWJn0KYktbBmi5yJUJuvnGuJkKhOAIEw6jjZ+6qdkZxdDAOtMNpK1+StPLEMQ
				lp2BuxOMMtz/uNt9CIMJ1Y9AlXGxmMYXb6J3Sf45P7Fiwh1jGsBpSR7JmdCrexPzLve2AUvGR2MC
				f5e/fgdlSUD6AftG9u4uuHklelSTI+gSluSn2EMvkVmZrs0ONXnSR4kG35bMurUPbkRjAYUhY/WI
				uJ0vAdIGOh9nBlUUkatg6RrzHvE/rayW80eJ5Q==
			</ds:SignatureValue>
		</ds:Signature>
		<saml:Subject xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
			<saml:NameID Format="urn:oasis:names:tc:SAML:2.0:nameid-format:transient" NameQualifier="https://siam.anoigo.nl/aselectserver/server" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
				5C1F9B3CCAAC0986B4D61AD04C098E90036D482915C79EC5121
				8E9B503F15A351922894D85B336D4F89FC65082C07AD18FC42CBFA356BD77A0C1B8059405ADE36EAD7D3531609
				D0611772F35ED5C0498B10A65FF8CC4AFEA64E0240D033584FBF6F001344BD1AD24688DC6A67EA75F1F4EBAB4E3
				BFE95D9
			</saml:NameID>
		</saml:Subject>
		<saml:AttributeStatement xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
			<saml:Attribute Name="uid" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
				<saml:AttributeValue xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="xs:string">
					1-900009676-Z-90000386-01.015-00000000
				</saml:AttributeValue>
			</saml:Attribute>
			<saml:Attribute Name="sel_level" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
				<saml:AttributeValue xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="xs:string">
					30
				</saml:AttributeValue>
			</saml:Attribute>
		</saml:AttributeStatement>
	</saml:Assertion>
	`)

	principal, err := decodeSAML(r)
	if err != nil {
		t.Fatalf("Error returned when decoding valid assertion: %v", err)
	}
	if principal.Uid != "1-900009676-Z-90000386-01.015-00000000" {
		t.Errorf("Incorrect principal in SAML assertion: %s", principal.Uid)
	}
	if principal.SelLevel != 30 {
		t.Errorf("Incorrect sel level in SAML assertion: %d", principal.SelLevel)
	}
}
