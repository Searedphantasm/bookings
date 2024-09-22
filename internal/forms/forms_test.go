package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Errorf("Form Valid() returned %v, wanted valid", isValid)
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields missing.")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest(http.MethodPost, "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows does not have required field when it does.")
	}

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has filed when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows min length for non-existing field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Should have an error but did not get one.")
	}

	postedData := url.Values{}
	postedData.Add("some_field", "some-value")
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("Form shows min length of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another_field", "some-value")
	form = New(postedData)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("Form shows min length of 1  is not met when it is")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("Should not have an error but got one.")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows valid email for non-existing field")
	}

	postedData = url.Values{}
	postedData.Add("email", "me@test.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Form shows invalid email for valid email")
	}

	postedData = url.Values{}
	postedData.Add("email", "x")
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Form shows valid email for invalid email")
	}

}
