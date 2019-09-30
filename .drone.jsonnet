local kobopatch = "0.14.0";
local versions = [
    "4.9.11311", "4.9.11314", "4.10.11591", "4.10.11655", "4.11.11911",
    "4.11.11976", "4.11.11980", "4.11.11982", "4.11.12019", "4.12.12111",
    "4.13.12638", "4.14.12777", "4.15.12920", "4.16.13162", "4.17.13651",
    "4.17.13694", "4.18.13737",
];

local pipeline(name, steps, opts) = std.mergePatch({kind: "pipeline", name: name, steps: steps}, opts);
local debian(name, commands) = {name: name, image: "debian:buster", commands: commands};
local debian_go(name, commands) = {name: name, image: "golang:1.13-buster", commands: commands};
local debian_pkgs(name, pkgs, commands) = {name: name, image: "debian:buster", commands: ["apt update", "apt install -y " + std.join(" ", pkgs)] + commands};
local plugin(name, plugin, settings) = {name: name, image: "plugins/" + plugin, settings: settings};
local depend_on(step, dependencies) = std.mergePatch(step, {depends_on: dependencies});

[
    pipeline("build", std.flattenArrays([
        [debian_go("build", ["go build -o ./scripts/build ./scripts/build.go", "./scripts/build -help || true"])],
        [depend_on(debian_go("download", ["./scripts/build -skipbuild"]), ["build"])],
        std.map((function(version) depend_on(debian_go(version, ["./scripts/build -skipdl " + version]), ["download"])), versions),
        [depend_on(debian("ls", ["ls -lah build"]), versions)],
    ]), {}),
    pipeline("test", std.flattenArrays([
        [debian_pkgs("kobopatch", ["wget", "ca-certificates"], [
            "wget -O kobopatch 'https://github.com/geek1011/kobopatch/releases/download/v" + kobopatch + "/kobopatch-linux-64bit'",
            "chmod +x kobopatch",
        ])],
        std.map((function(version) depend_on(debian(version, [
            'bash -c \'export PATH="$PWD:$PATH"; set -o pipefail; ./scripts/test.sh ' + version + ' | tr "\\r" "\\n"\'',
        ]), ["kobopatch"])), versions)
    ]), {}),
    pipeline("build-release", std.flattenArrays([
        [debian_go("build", ["go build -o ./scripts/build ./scripts/build.go", "./scripts/build -help || true"])],
        [depend_on(debian_go("download", ["./scripts/build -skipbuild"]), ["build"])],
        std.map((function(version) depend_on(debian(version, ["./scripts/build -skipdl " + version]), ["download"])), versions),
        [depend_on(plugin("release", "github-release", {
            api_key: {from_secret: "github_token"},
            files: ["build/*"],
            title: "${DRONE_TAG}",
        }), versions)],
    ]), {
        trigger: {ref: ["refs/tags/v*"]},
        depends_on: ["test"],
    }),
]
