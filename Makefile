lint:
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 --sort-results

# XXX careful with this if you are using workspaces - it will trash stuff
lint-fix:
	golangci-lint run --fix