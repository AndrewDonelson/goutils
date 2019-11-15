package goutils

import (
	"strings"
)

// HTML5_Page returns string content representing a valid HTML5 page with the given title set and the content provided placed in a
// DIV an ID  of content
func HTML5_Page(title, content string) string {
	return strings.Join([]string{"<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\" />\n<title>", title, "</title>\n<link rel=\"shortcut icon\" href=\"http://52.3.210.246/wp-content/uploads/2015/09/favicon.ico\" type=\"image/x-icon\">\n<style>\nhtml,body {height:100%;width:100%;margin:0;}\nbody , body {display:flex;}\n#content,form {margin:auto;}\n</style>\n</head>\n<body>\n<div id=\"content\">", content, "</div>\n</body>\n</html>\n"}, "")
}

// HTML5_Form_Login returns string content representing a basic login form centered on the page
func HTML5_Form_Login() string {
	return "<form id=\"form_login\" action=\"/v0/user/login\" method=\"post\">\n<h1>Test Login Page</h1><p><input type=\"text\" id=\"email\" name=\"email\" required placeholder=\"account email\" /></p>\n<p><input type=\"password\" id=\"password\" name=\"password\" required placeholder=\"password\" /></p>\n<p><button id=\"submitbutton\" type=\"submit\">Login</button></p>\n</form>\n"
}

// HTML5_Page_Login returns a complete basic Login page
func HTML5_Page_Login() string {
	return HTML5_Page("Member Login", HTML5_Form_Login())
}

func HTML5_Page_Not_Implemented(name string) string {
	return HTML5_Page("Forgot Password", "<h1>"+name+" Not yet implemented</h1>")
}
