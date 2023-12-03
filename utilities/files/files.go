package files

import (
	"fmt"
	"os"
)

func GetPermissions(fileInfo os.FileInfo) string {
	mode := fileInfo.Mode()
	permissions := ""

	// Owner
	permissions += getPermissionChar(mode, os.ModePerm&0400, 'r')
	permissions += getPermissionChar(mode, os.ModePerm&0200, 'w')
	permissions += getPermissionChar(mode, os.ModePerm&0100, 'x')

	// Group
	permissions += getPermissionChar(mode, os.ModePerm&0040, 'r')
	permissions += getPermissionChar(mode, os.ModePerm&0020, 'w')
	permissions += getPermissionChar(mode, os.ModePerm&0010, 'x')

	// Others
	permissions += getPermissionChar(mode, os.ModePerm&0004, 'r')
	permissions += getPermissionChar(mode, os.ModePerm&0002, 'w')
	permissions += getPermissionChar(mode, os.ModePerm&0001, 'x')

	return permissions
}

func getPermissionChar(mode os.FileMode, mask os.FileMode, char rune) string {
	if mode&mask != 0 {
		return string(char)
	}
	return "-"
}

func FormatSize(size int64, fileType string) string {
	if fileType == "Directory" {
		return ""
	}
	
	const (
		KB = 1 << 10
		MB = 1 << 20
	)

	switch {
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}
