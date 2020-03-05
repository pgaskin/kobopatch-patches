#!/bin/bash
cd "$(dirname "$0")"
rm -f "out/KoboRoot.tgz"
case `uname -s` in
    CYGWIN_NT*)
        ./bin/koboptch-windows.exe
        ;;
    Darwin)
	    ./bin/kobopatch-darwin-64bit
	    ;;
    Linux)
	    case `uname -m` in
	        i?86)
                ./bin/kobopatch-linux-32bit
		        ;;
	        x86_64)
                ./bin/kobopatch-linux-64bit
		        ;;
            arm*)
                ./bin/kobopatch-linux-arm
                ;;
            aarch64)
                ./bin/kobopatch-linux-arm
                ;;
            *)
                echo "Unsupported architecture"
	esac
	;;
esac