
fmt:
	find . -name '*.go' -exec gofumpt -w -s -extra {} \;

doc:
	find . -name '*.go' -exec code2prompt --template ~/code2prompt/templates/document-the-code.hbs --output {}.md {} \;