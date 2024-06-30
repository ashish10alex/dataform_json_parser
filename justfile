build:
    go build -o bin/dj
    cp bin/dj ~/bin/dj

generate_completion:
    bin/dj completion zsh > ~/.zsh/completion/dj

bbuild: build generate_completion

test_table_ops:
    go run . --json-file test/op.json table-ops cost --all

test_tag_ops:
    go run . --json-file test/op.json tag-ops cost --all

test_single_table_ops_cost:
    go run . --json-file test/dpr.json table-ops cost --file=0400_WIP_HISTORY

test_compiled_query:
    go run . --json-file test/dpr.json table-ops query --file=0400_WIP_HISTORY

debug:
    dlv debug -- -j test/dpr.json tag-ops -u

download_a_release_from_github_mac_arm:
    curl -OL https://github.com/ashish10alex/dj/releases/download/v0.0.6-pre/dj_Darwin_arm64.tar.gz ; tar -xzf dj_Darwin_arm64.tar.gz

download_a_release_from_github_linux:
    curl -OL https://github.com/ashish10alex/dj/releases/download/v0.0.6-pre/dj_Linux_x86_64.tar.gz ; tar -xzf dj_Linux_x86_64.tar.gz

clean_cache:
    go clean --cache

test:
    gotestsum --format testname
