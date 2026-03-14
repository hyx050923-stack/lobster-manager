package platform

import "runtime"

type Platform string

const (
	Linux   Platform = "linux"
	Windows Platform = "windows"
	Android Platform = "android"
)

func Detect() Platform {

	if runtime.GOOS == "android" {
		return Android
	}

	if runtime.GOOS == "windows" {
		return Windows
	}

	return Linux
}