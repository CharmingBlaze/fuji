# dist-template

This folder defines the contents of the release SDK zip.

When `scripts/build-release.sh` (or `.ps1`) runs, it copies:
- The `fuji` binary
- The `fujiwrap` binary  
- Everything in this `dist-template/` folder

Into a zip named `fuji-vX.Y.Z-sdk-<platform>.zip`.

Do not put compiler internals, test files, or build scripts in this folder.
Users will see everything here.
