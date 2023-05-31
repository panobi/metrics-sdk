package panobi

func errorIs(want string, got error) bool {
	if got == nil {
		return want == ""
	} else {
		return want == got.Error()
	}
}
