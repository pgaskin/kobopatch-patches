local kobopatch = "0.13.0";

local pipeline(name, steps) = {kind: "pipeline", name: name, steps: steps};
local debian(name, commands) = {name: name, image: "debian:stretch", commands: commands};

local test(version) = {
    name: version,
    image: "debian:stretch",
    commands: [
        "export PATH=\"$PWD:$PATH\"",
        "./scripts/test.sh " + version,
    ],
    depends_on: [
        "kobopatch"
    ],
};

[
    pipeline("build", [
        debian("build", [
            "apt update",
            "apt install -y dos2unix wget zip ca-certificates",
            "./scripts/build.sh"
        ]),
        debian("release", [
            "echo 'Not Implemented'",
        ]),
    ]),
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
]