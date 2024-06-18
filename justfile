build:
    go build -o bin/dj
    cp bin/dj ~/bin/dj

test:
    bin/dj -h

generate_completion:
    bin/dj completion zsh > ~/.zsh/completion/dj

bbuild: build generate_completion

test_table_ops:
    go run . --json-file test/op.json table-ops cost --all

test_tag_ops:
    go run . --json-file test/op.json tag-ops cost --all

test_single_table_ops_cost:
    go run . --json-file test/dpr.json table-ops cost --table=table_name

debug:
    dlv debug -- -j test/dpr.json tag-ops -u
