package platform

import "runtime"

func IsAndroid() bool {
	return runtime.GOOS == "android"
}