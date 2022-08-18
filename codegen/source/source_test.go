package source

import "testing"

func TestNewUserDescriptionParser(t *testing.T) {
	type expectation struct {
		give string
		want string
	}
	cases := []struct {
		regex        string
		words        []string
		expectations []expectation
	}{
		{
			regex: `[$]+`,
			words: []string{"remove", "test phrase"},
			expectations: []expectation{
				{give: "my $$phrase", want: "my phrase"},
				{give: "my test phrase", want: "my"},
				{give: "remove my test phrase", want: "my"},
				{give: "remove link: http://cloudquery.io", want: "link: http://cloudquery.io"},
				{give: "first sentence. second sentence.", want: "first sentence. second sentence."},
			},
		},
	}
	for _, tc := range cases {
		p := NewUserDescriptionParser(tc.regex, tc.words)
		for _, exp := range tc.expectations {
			got := p.Parse(exp.give)
			if got != exp.want {
				t.Errorf("regex=%q\nwords=%q\nParse(%q) =\n     %q,\nwant %q",
					tc.regex, tc.words, exp.give, got, exp.want)
			}
		}
	}
}
