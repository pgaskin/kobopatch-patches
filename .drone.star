kobopatch = "0.15.0"
versions = [
    "4.10.11591", "4.10.11655", "4.11.11911", "4.11.11976", "4.11.11980",
    "4.11.11982", "4.11.12019", "4.12.12111", "4.13.12638", "4.14.12777",
    "4.15.12920", "4.16.13162", "4.17.13651", "4.17.13694", "4.18.13737",
    "4.19.14123", "4.20.14601", "4.20.14617", "4.20.14622", "4.21.15015",
    "4.22.15190", "4.22.15268", "4.23.15505", "4.23.15548", "4.24.15672",
    "4.24.15676", "4.25.15875", "4.26.16704", "4.28.17623", "4.28.17820",
    "4.28.17826", "4.28.17925", "4.28.18220",
]

pipeline = [{
    "kind": "pipeline",
    "name": "build",
    "steps": [{
        "name": "build",
        "image": "golang:1.13-buster",
        "commands": [
            "go build -o ./scripts/buildscript ./scripts/build",
            "./scripts/buildscript -help || true",
        ],
    }, {
        "name": "download",
        "image": "golang:1.13-buster",
        "commands": ["./scripts/buildscript -kpver v%s -skipbuild" % kobopatch],
        "depends_on": ["build"],
    }] + [{
        "name": version,
        "image": "debian:buster",
        "commands": ["./scripts/buildscript -kpver v%s -skipdl %s" % (kobopatch, version)],
        "depends_on": ["download"],
    } for version in versions] + [{
        "name": "ls",
        "image": "debian:buster",
        "commands": ["ls -lah build"],
        "depends_on": versions,
    }],
}, {
    "kind": "pipeline",
    "name": "test",
    "steps": [{
        "name": "testdata",
        "image": "golang:1.13-buster",
        "commands": [
            "git fetch --depth=1 https://github.com/pgaskin/kobopatch-testdata v1",
            "git worktree add testdata FETCH_HEAD",
        ],
    }, {
        "name": "kobopatch",
        "image": "debian:buster",
        "commands": [
            "apt update",
            "apt install -y wget ca-certificates",
            "wget -O kobopatch 'https://github.com/pgaskin/kobopatch/releases/download/v%s/kobopatch-linux-64bit'" % kobopatch,
            "chmod +x kobopatch",
        ],
    }, {
        "name": "build",
        "image": "golang:1.13-buster",
        "commands": [
            "go build -o ./scripts/testscript ./scripts/test",
            "./scripts/testscript -help || true",
        ],
    }] + [{
        "name": version,
        "image": "debian:buster",
        "commands": ["./scripts/testscript -kpbin ./kobopatch %s" % version],
        "depends_on": ["kobopatch", "build", "testdata"],
    } for version in versions],
}, {
    "kind": "pipeline",
    "name": "build-release",
    "steps": [{
        "name": "build",
        "image": "golang:1.13-buster",
        "commands": [
            "go build -o ./scripts/buildscript ./scripts/build",
            "./scripts/buildscript -help || true",
        ],
    }, {
        "name": "download",
        "image": "golang:1.13-buster",
        "commands": ["./scripts/buildscript -kpver v%s -skipbuild" % kobopatch],
        "depends_on": ["build"],
    }] + [{
        "name": version,
        "image": "debian:buster",
        "commands": ["./scripts/buildscript -kpver v%s -skipdl %s" % (kobopatch, version)],
        "depends_on": ["download"],
    } for version in versions] + [{
        "name": "release",
        "image": "plugins/github-release",
        "settings": {
            "api_key": {"from_secret": "github_token"},
            "files": ["build/*"],
            "title": "${DRONE_TAG}",
        },
        "depends_on": versions,
    }],
    "trigger": {"ref": ["refs/tags/v*"]},
    "depends_on": ["test"],
}]

def main(ctx):
  return pipeline
