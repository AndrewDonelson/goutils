package goutils

import "testing"

func TestHTML5Page(t *testing.T) {
	got := HTML5Page("Test Title", "<p>Example Paragraph</p>")
	//println(len(got))
	Equals(t, len(got), 403)
}

func TestHTML5FormLogin(t *testing.T) {
	got := HTML5FormLogin()
	Equals(t, len(got), 340)
}

func TestHTML5PageLogin(t *testing.T) {
	got := HTML5PageLogin()
	Equals(t, len(got), 721)
}

func TestPageNotImplemented(t *testing.T) {
	got := HTML5PageNotImplemented("Some Page")
	Equals(t, len(got), 421)
}
