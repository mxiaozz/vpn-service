package util

import (
	"fmt"
	"strings"
)

func If[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}

func IsAdminId(userId int64) bool {
	return userId == 1
}

func IsHttp(link string) bool {
	link = strings.TrimLeft(link, " ")
	if len(link) > 7 {
		link = link[0:7]
	}

	list := []string{"http://", "https://"}
	for _, v := range list {
		if strings.EqualFold(link, v) {
			return true
		}
	}

	return false
}

func HumanByteSize(dataSize int64) string {
	var kb int64 = 1024
	var mb int64 = 1024 * kb
	var gb int64 = 1024 * mb
	var tb int64 = 1024 * gb

	size := float64(0)
	if dataSize > tb {
		size = float64(dataSize) / float64(tb)
		return fmt.Sprintf("%.2f TB", size)
	}
	if dataSize > gb {
		size = float64(dataSize) / float64(gb)
		return fmt.Sprintf("%.2f GB", size)
	}
	if dataSize > mb {
		size = float64(dataSize) / float64(mb)
		return fmt.Sprintf("%.2f MB", size)
	}
	if dataSize > kb {
		size = float64(dataSize) / float64(kb)
		return fmt.Sprintf("%.2f KB", size)
	}
	return fmt.Sprintf("%d B", dataSize)
}
