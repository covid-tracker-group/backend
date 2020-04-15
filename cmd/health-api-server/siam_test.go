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

	r.Header.Add("X-saml-attribute-token1", "cGFydDAx")
	if string(getSAMLAssertionFromHTTPHeader(r)) != "part01" {
		t.Fatal("Failed to find single-header SAML Assertion from HTTP headers")
	}

	r.Header.Add("X-saml-attribute-token2", "IGFuZCBzZW")
	r.Header.Add("X-saml-attribute-token3", "NvbmQgb25l")
	assertion := string(getSAMLAssertionFromHTTPHeader(r))
	if assertion != "part01 and second one" {
		t.Fatalf("Data from consecutive headers not read %v", assertion)
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

	// This token was created using a test UZI card.
	r.Header.Add("X-saml-attribute-token1",
		`PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHNhbWwyOkFzc2VydGlvbiBJRD0iSUM1MzJBQ0U0NDMxOENGN0Q0OUM1RDI3QjBCMDc1RUQ0QzIzODA0NEQiIElzc3VlSW5zdGFudD0iMjAyMC0wNC0xNFQxODoyNjowOC4wNTFaIiBWZXJzaW9uPSIyLjAiIHhtbG5zOnNhbWwyPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIiB4bWxuczp4cz0iaHR0cDovL3d3dy53My5vcmcvMjAwMS9YTUxTY2hlbWEiPjxzYW1sMjpJc3N1ZXIgRm9ybWF0PSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bmFtZWlkLWZvcm1hdDplbnRpdHkiIHhtbG5zOnNhbWwyPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIj5odHRwczovL3NpYW0xLnRlc3QuYW5vaWdvLm5sL2FzZWxlY3RzZXJ2ZXIvc2VydmVyPC9zYW1sMjpJc3N1ZXI+PGRzOlNpZ25hdHVyZSB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+PGRzOlNpZ25lZEluZm8+PGRzOkNhbm9uaWNhbGl6YXRpb25NZXRob2QgQWxnb3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiLz48ZHM6U2lnbmF0dXJlTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnI3JzYS1zaGExIi8+PGRzOlJlZmVyZW5jZSBVUkk9IiNJQzUzMkFDRTQ0MzE4Q0Y3RDQ5QzVEMjdCMEIwNzVFRDRDMjM4MDQ0RCI+PGRzOlRyYW5zZm9ybXM+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNlbnZlbG9wZWQtc2lnbmF0dXJlIi8+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMTAveG1sLWV4Yy1jMTRuIyI+PGVjOkluY2x1c2l2ZU5hbWVzcGFjZXMgUHJlZml4TGlzdD0ieHMiIHhtbG5zOmVjPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiLz48L2RzOlRyYW5zZm9ybT48L2RzOlRyYW5zZm9ybXM+PGRzOkRpZ2VzdE1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNzaGExIi8+PGRzOkRpZ2VzdFZhbHVlPnNVRGJDaUJSYlZnTDFCOE5tazE4OTVUU2pLbz08L2RzOkRpZ2VzdFZhbHVlPjwvZHM6UmVmZXJlbmNlPjwvZHM6U2lnbmVkSW5mbz48ZHM6U2lnbmF0dXJlVmFsdWU+Qmdzc3JobmJHQnhpYWVSOHUzQmh5SEp4WUNlMjV0QXpqaE8zVi8xODRFWGljL0xHcGdaTVQzUmp4OGtrOXFWT1QvV1owNnpleDlreWI3dWY2VURjZnlhcVhOZEdDUWRMQ21KQ0U1Ukc3SVZWclBOeWROUTM2RnZuVzl6Zkw1NVR4RklFcTJuckpVeVJIc1RWbGJvTGpsQnpCQlRaYVpaS1c5WFhnQ2tHRlRYc2VaMWZOQ0VoZWNXSDl5U1QxNGx4T0pGeTA2WlNSMkdRRHY4d2s2Z243UWZBTDJsUHViQ2xxbndQakdydjI3TVUzeTdJNlhTWlZUcHZuQzdSNWl0WUtaTGwxKzdKaWdBNTRnQ0VRUThjcGdpZDlBREk3WkhIRXhUdHlWYmdOVTh1dGRxSkg4d3BBZk8w`)
	r.Header.Add("X-saml-attribute-token2",
		`Vm44WmRjVmsrdXExQVpLTFVXbENXMGM4YnJTOEx3PT08L2RzOlNpZ25hdHVyZVZhbHVlPjwvZHM6U2lnbmF0dXJlPjxzYW1sMjpTdWJqZWN0IHhtbG5zOnNhbWwyPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIj48c2FtbDI6TmFtZUlEIEZvcm1hdD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm5hbWVpZC1mb3JtYXQ6dHJhbnNpZW50IiBOYW1lUXVhbGlmaWVyPSJodHRwczovL3NpYW0xLnRlc3QuYW5vaWdvLm5sL2FzZWxlY3RzZXJ2ZXIvc2VydmVyIj4wM0EzNkQxODZGRjAzRDI1MUQ1NDRGMTgxQUVCQzgyNkJENkM1QTQ4MTc1MzVBNjdENzBBRUYzNjEzNENFODcwQkZBODgxRjQwMjQ5MkNFQTc1OUQ3MkFDQjQwOUI3Q0EyNEJEMjM1RjZDQUIyMkVFNzdCMEZFMjNFMTYzREIxRTI1Rjc0MzBBNjY4RjhDQzIxOEU3ODM4RDI3RDEzRTc2OUY1QzgxM0U5NTRDMjJBQTNGNTgxOTZGREExOTRENDdGOTIzRjJGNDNFMjg0REFFOUU1NDYwQzlGNEREMTY3MjRBMkI1MUEyODhEMkZFQ0Y8L3NhbWwyOk5hbWVJRD48L3NhbWwyOlN1YmplY3Q+PHNhbWwyOkF0dHJpYnV0ZVN0YXRlbWVudCB4bWxuczpzYW1sMj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmFzc2VydGlvbiI+PHNhbWwyOkF0dHJpYnV0ZSBOYW1lPSJ1aWQiPjxzYW1sMjpBdHRyaWJ1dGVWYWx1ZSB4bWxuczp4cz0iaHR0cDovL3d3dy53My5vcmcvMjAwMS9YTUxTY2hlbWEiIHhtbG5zOnhzaT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS9YTUxTY2hlbWEtaW5zdGFuY2UiIHhzaTp0eXBlPSJ4czpzdHJpbmciPjEtOTAwMDE4MTAxLVotOTAwMDAzODEtMDEuMDE1LTAwMDAwMDAwPC9zYW1sMjpBdHRyaWJ1dGVWYWx1ZT48L3NhbWwyOkF0dHJpYnV0ZT48c2FtbDI6QXR0cmlidXRlIE5hbWU9Imxhc3RzeW5jdGltZSI+PHNhbWwyOkF0dHJpYnV0ZVZhbHVlIHhtbG5zOnhzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYSIgeG1sbnM6eHNpPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYS1pbnN0YW5jZSIgeHNpOnR5cGU9InhzOnN0cmluZyI+MTU4Njg4ODc2Nzk5NDwvc2FtbDI6QXR0cmlidXRlVmFsdWU+PC9zYW1sMjpBdHRyaWJ1dGU+PC9zYW1sMjpBdHRyaWJ1dGVTdGF0ZW1lbnQ+PC9zYW1sMjpBc3NlcnRpb24+`)

	raw_assertion := getSAMLAssertionFromHTTPHeader(r)
	if len(raw_assertion) != 2661 {
		t. Fatalf("Got unexpected length while extracting assertion from request: %d", len(raw_assertion))
	}
	principal, err := decodeSAML(r)
	if err != nil {
		t.Fatalf("Error returned when decoding valid assertion: %v", err)
	}
	if principal.Uid != "1-900018101-Z-90000381-01.015-00000000" {
		t.Errorf("Incorrect principal in SAML assertion: %s", principal.Uid)
	}
	if principal.SelLevel != 0 {
		t.Errorf("Incorrect sel level in SAML assertion: %d", principal.SelLevel)
	}
}
