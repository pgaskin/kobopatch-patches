local kobopatch = "0.13.0";

local pipeline(name, steps) = {kind: "pipeline", name: name, steps: steps};
local debian(name, commands) = {name: name, image: "debian:stretch", commands: commands};
local golang(name, commands) = {name: name, image: "golang:1.12-stretch", commands: commands};

local test(version) =
    std.mergePatch(debian(version, [
        "export PATH=\"$PWD:$PATH\"",
        "./scripts/test.sh " + version + ' | tr "\\r" "\\n"',
    ]), {
        depends_on: ["kobopatch"],
    });

[
    std.mergePatch(pipeline("build", [
        debian("build", [
            "apt update",
            "apt install -y dos2unix wget zip ca-certificates",
            "./scripts/build.sh"
        ]),
    ]), {
        trigger: {ref: {exclude: ["refs/tags/v*"]}},
    }),
    pipeline("test", [
        debian("kobopatch", [
            "apt update",
            "apt install -y wget ca-certificates",
            "wget -O kobopatch 'https://github.com/geek1011/kobopatch/releases/download/v" + kobopatch + "/kobopatch-linux-64bit'",
            "chmod +x kobopatch",
        ]),
        test("4.9.11311"),
        test("4.9.11314"),
        test("4.10.11591"),
        test("4.10.11655"),
        test("4.11.11911"),
        test("4.11.11976"),
        test("4.11.11980"),
        test("4.11.11982"),
        test("4.11.12019"),
        test("4.12.12111"),
        test("4.13.12638"),
        test("4.14.12777"),
        test("4.15.12920"),
    ]),
    std.mergePatch(pipeline("build-release", [
        debian("build", [
            "apt update",
            "apt install -y dos2unix wget zip ca-certificates",
            "./scripts/build.sh"
        ]),
        golang("release", [
            "GO111MODULE=on go get github.com/tcnksm/ghr@v0.12.1",
            "ghr ${DRONE_TAG} build/",
        ]),
    ]), {
        trigger: {ref: ["refs/tags/v*"]},
        depends_on: ["test"],
    }),
]