local kobopatch = "0.14.0";
local versions = [
    "4.9.11311", "4.9.11314", "4.10.11591", "4.10.11655", "4.11.11911",
    "4.11.11976", "4.11.11980", "4.11.11982", "4.11.12019", "4.12.12111",
    "4.13.12638", "4.14.12777", "4.15.12920", "4.16.13162",
];

local pipeline(name, steps, opts) = std.mergePatch({kind: "pipeline", name: name, steps: steps}, opts);
local debian(name, commands) = {name: name, image: "debian:buster", commands: commands};
local debian_pkgs(name, pkgs, commands) = {name: name, image: "debian:buster", commands: ["apt update", "apt install -y " + std.join(" ", pkgs)] + commands};
local plugin(name, plugin, settings) = {name: name, image: "plugins/" + name, settings: settings};

[
    pipeline("build", [
        debian_pkgs("build", ["dos2unix", "wget", "zip", "ca-certificates"], ["./scripts/build.sh"]),
    ], {
        trigger: {ref: {exclude: ["refs/tags/v*"]}},
    }),
    pipeline("test", std.flattenArrays([
        [debian_pkgs("kobopatch", ["wget", "ca-certificates"], [
            "wget -O kobopatch 'https://github.com/geek1011/kobopatch/releases/download/v" + kobopatch + "/kobopatch-linux-64bit'",
            "chmod +x kobopatch",
        ])],
        std.map((function(version) std.mergePatch(debian(version, [
            'export PATH="$PWD:$PATH"',
            'set -o pipefail',
            './scripts/test.sh ' + version + ' | tr "\\r" "\\n"',
        ]), {
            depends_on: ["kobopatch"],
        })), versions)
    ]), {}),
    pipeline("build-release", [
        debian_pkgs("build", ["dos2unix", "wget", "zip", "ca-certificates"], ["./scripts/build.sh"]),
        plugin("release", "github-release", {
            api_key: {from_secret: "github_token"},
            files: ["build/*"],
            title: "${DRONE_TAG}",
        }),
    ], {
        trigger: {ref: ["refs/tags/v*"]},
        depends_on: ["test"],
    }),
]