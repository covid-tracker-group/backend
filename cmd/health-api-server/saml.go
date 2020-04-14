package main

import (
	"encoding/xml"
)

type SAMLIssuer struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	Format  string   `xml:"Format,attr"`
	Name    string   `xml:",chardata"`
}

type XMLDigMethod struct {
	Value string `xml:"Algorithm,attr"`
}

type XMLDSigReference struct {
	XMLName      xml.Name       `xml:"http://www.w3.org/2000/09/xmldsig# Reference"`
	Transforms   []XMLDigMethod `xml:"Transforms>Transform"`
	DigestMethod XMLDigMethod   `xml:"DigestMethod"`
	DigestValue  string         `xml:"DigestValue"`
}

type XMLDigSignedInfo struct {
	XMLName                xml.Name         `xml:"http://www.w3.org/2000/09/xmldsig# SignedInfo"`
	CanonicalizationMethod XMLDigMethod     `xml:"CanonicalizationMethod"`
	SignatureMethod        XMLDigMethod     `xml:"SignatureMethod"`
	Reference              XMLDSigReference `xml:"Reference"`
}

type XMLDSigSignature struct {
	XMLName        xml.Name         `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
	SignatureValue string           `xml:"SignatureValue"`
	SignedInfo     XMLDigSignedInfo `xml:"SignedInfo"`
}

type SAMLNameID struct {
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion NameID"`
	Format        string   `xml:"Format,attr"`
	NameQualifier string   `xml:"NameQualifier,attr"`
	Value         string   `xml:",chardata"`
}

type SAMLAttributeValue struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeValue"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:",chardata"`
}
type SAMLAttribute struct {
	XMLName xml.Name             `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute"`
	Name    string               `xml:"Name,attr"`
	Values  []SAMLAttributeValue ` xml:"AttributeValue"`
}

type SAMLAssertion struct {
	XMLName      xml.Name         `xml:"urn:oasis:names:tc:SAML:2.0:assertion Assertion"`
	Version      string           `xml:"Version,attr"`
	IssueInstant string           `xml:"IssueInstant,attr"`
	Issuer       SAMLIssuer       `xml:"Issuer"`
	Signature    XMLDSigSignature `xml:"Signature"`
	Subject      SAMLNameID       `xml:"Subject>NameID"`
	Attributes   []SAMLAttribute  `xml:"AttributeStatement>Attribute"`
}
