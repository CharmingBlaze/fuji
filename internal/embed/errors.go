package fujembed

import "errors"

// ErrDevelopmentBuild means the binary was built without the release build tag; no assets are embedded.
var ErrDevelopmentBuild = errors.New("not a release build")
