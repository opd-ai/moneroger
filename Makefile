
fmt:
	find . -name '*.go' -exec gofumpt -w -s -extra {} \;

doc:
	find ./*/ -type d -exec code2prompt --template ~/code2prompt/templates/write-a-test.hbs --output {}/tests.md {} \;

readme:
	code2prompt --template ~/code2prompt/templates/write-github-readme.hbs --output readme_.md .