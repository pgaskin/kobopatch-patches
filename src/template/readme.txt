This is for use with firmware {{version}}.

=== How to use ===
Make sure your kobo is already running the firmware version you are trying to patch.

1. Download the firmware from https://pgaskin.net/KoboStuff/kobofirmware.html to the src folder.
   The zip should be called something like kobo-update-{{version}}.zip. You may need to rename it.
2. Enable patches in the files in the src folder (or use the overrides in kobopatch.yaml).
3. Run kobopatch.bat on Windows, or kobopatch.sh on Linux.
4. If the patching succeeded, a file named KoboRoot.tgz will be created.
   Copy it to the .kobo folder of your device.

=== To update kobopatch ===
1. Download the updated version from https://github.com/pgaskin/kobopatch/releases

=== To revert the patches ===
1. Disable all patches.
2. Repeat steps 3-4 above.

You can also reinstall the firmware.

=== Translations ===
kobopatch also supports compiling and adding ts files. lrelease should be in your PATH, or the path
should be specified in kobopatch.yaml.

=== To report bugs: ===
Open an issue on https://github.com/pgaskin/kobopatch/issues, or respond to the thread on MobileRead.
Make sure you provide log.txt