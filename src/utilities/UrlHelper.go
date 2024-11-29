package utilities

import (
	"fmt"
	"strings"
)

// UrlJoin is a utility function for constructing a URL by joining its components.
//
// Parameters:
//   - protocol (string): The URL protocol (e.g., "http" or "https").
//   - host (string): The URL host (e.g., "example.com" or "localhost").
//   - port (int): The port number to include in the URL.
//   - paths (string variadic): One or more path components to be joined in the URL. These can include
//     directories and file names, and they will be cleaned up and properly formatted in the URL.
//
// Returns:
//   - string: The fully constructed URL as a string.
//
// Steps:
//   1. Initialize an empty string called 'builder' to store the joined path components.
//   2. Determine the number of path components in the 'paths' variadic input.
//   3. Iterate over the 'paths' and add each component to the 'builder' after trimming any leading or
//      trailing slashes.
//   4. If the current path is not the last one (index < pathSize-1), append a forward slash '/' to
//      separate path components.
//   5. Format and construct the complete URL using the 'fmt.Sprintf' function, including the 'protocol',
//      'host', 'port', and the 'builder' containing joined path components.
//   6. Return the fully constructed URL as a string.

func UrlJoin(protocol string, host string, port int, paths ...string) string {
	builder := ""
	pathSize := len(paths)
	for index, p := range paths {
		builder += strings.Trim(p, "/")
		if index < pathSize-1 {
			builder += "/"
		}
	}
	return fmt.Sprintf("%s://%s:%d/%s", protocol, host, port, builder)
}
