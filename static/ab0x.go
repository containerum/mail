// Code generated by fileb0x at "2018-06-18 10:13:20.141018254 +0300 MSK m=+0.020356195" from config file "b0x.yaml" DO NOT EDIT.
// modification hash(5fbbb523fd3849ef057c2a09ad772c1b.3daf2b4d0cc8675651cd8e10adf6cbc7)

package static

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path"

	"context"
	"golang.org/x/net/webdav"
)

var (
	// CTX is a context for webdav vfs
	CTX = context.Background()

	// FS is a virtual memory file system
	FS = webdav.NewMemFS()

	// Handler is used to server files through a http handler
	Handler *webdav.Handler

	// HTTP is the http file system
	HTTP http.FileSystem = new(HTTPFS)
)

// HTTPFS implements http.FileSystem
type HTTPFS struct{}

// FileSwaggerJSON is "swagger.json"
var FileSwaggerJSON = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xec\x5c\xcd\x6e\xdc\x38\x12\xbe\xe7\x29\x08\xcd\x02\x3e\xac\xbb\x95\x0d\x16\x7b\xf0\x6d\xb0\xed\x49\x8c\x99\x59\x78\x6d\x27\x97\xb1\x61\xd0\x62\xb5\x9a\x89\x44\xca\x24\xe5\xd8\x63\xf8\xdd\x07\x24\xf5\x2f\xaa\xa5\xfe\x53\xbb\x33\xce\x25\x6a\xa9\x58\x24\xab\xbe\xaf\xaa\x44\x52\x7e\x7e\x87\x90\x17\x70\x26\xd3\x18\xa4\x77\x82\xfe\x78\x87\x10\x42\x1e\x4e\x92\x88\x06\x58\x51\xce\xfc\xaf\x92\x33\xef\x1d\x42\x37\xc7\x5a\x36\x11\x9c\xa4\xc1\x30\x59\xf9\x1d\x87\x21\x08\xef\x04\x79\x1f\xa6\xef\x3d\x73\x8f\xb2\x39\xf7\x4e\xd0\xb3\x6d\x4b\x40\x06\x82\x26\xba\xad\x96\xfa\x1d\xd3\x68\x72\x05\x71\x12\x61\x05\x02\x49\x10\x0f\x34\x00\x44\x65\x71\x39\xe7\xfa\x36\x23\x94\x85\x08\x62\x4c\x23\x39\x35\x7a\x11\xf2\x14\x55\x11\xb4\xb5\xe8\xe1\xbc\xd8\xa1\x63\xb5\x90\x65\xdf\x7e\x0c\x52\xe2\x10\xca\x5b\x08\x79\x21\xa8\xca\xcf\xf6\x08\x17\x4a\x25\xf2\xc4\xf7\x83\xc5\x34\xd1\x6d\xa7\x01\x67\x0a\x53\x06\x22\x8d\xa7\x0c\x94\x8f\x13\x3a\x21\x3c\x90\x7e\xcc\x49\x1a\x81\xf4\x83\xc5\x44\x8f\x73\xa2\xb2\x01\xf9\x94\x11\x78\x9c\x2e\x54\x1c\xfd\x14\x82\x9a\x14\xa3\x38\x2e\x3b\x55\x38\x2c\x0d\x9c\xdd\xfb\x3d\x97\x2b\x6e\xde\x54\x5a\xc8\x34\x8e\xb1\x78\xd2\x43\xfc\x08\x0a\xe5\x4a\x51\x44\xa5\x9a\x56\x55\xf3\x04\x84\xf1\xd5\x19\x31\xb6\xb2\x82\xbf\x51\xa9\x3e\x82\xfa\x84\x19\x89\x40\x54\xe5\x13\x2c\x70\x0c\x0a\x44\x73\x40\xcf\x95\x6b\x84\xbc\x7f\x08\x98\x6b\x85\x3f\xf9\x65\x0b\xff\xb3\x04\x71\xc1\x23\xf8\x04\x98\x58\x4f\xe4\xff\x5e\x8e\xbb\x55\xa9\xa7\xc4\xf8\x51\x2a\x41\x59\xe8\x1d\xd7\x9f\x32\x1c\x9b\xa7\xda\xfc\xcd\x67\xd4\xf8\xe8\x3e\x05\xf1\xb4\xe5\xce\x40\xdc\xae\xd2\xa1\xd3\x47\x02\x64\xc2\x99\xac\x01\xce\x3c\xf8\xf0\xfe\x7d\xe3\x56\x1b\x78\x99\x47\x8d\x43\x51\x08\x0a\xe5\xda\x9a\x43\x92\xc1\x02\x62\xdc\xd2\x57\x73\x11\x81\x39\x65\x54\xab\x96\x7e\x05\x01\x17\xb9\xca\x5a\xcb\x97\x2e\x4b\x7a\x04\xe6\x38\x8d\x54\x7b\xec\x65\x4f\xc5\x9c\x7d\x10\x82\x0b\xb7\x91\x2a\x5a\xbd\xc7\x49\x0c\x6a\xc1\xc9\xe4\x81\x4a\x7a\x47\x23\xaa\x0c\xa8\x13\x41\x1f\xb0\x2a\x46\x66\xdb\x66\xed\x4a\x26\xfb\xcf\xd9\xd5\x2d\x25\x2f\xaf\x82\xd6\x93\x80\x27\x4f\x5b\xa5\xb6\xa4\x2c\x8c\x20\x67\xf8\x00\x6e\x1f\x18\xaf\x4b\x0f\xba\xc9\xa6\xa3\x78\xf3\x89\x80\xfb\x94\x0a\xd0\xb3\x56\x22\x85\x71\x98\xb8\x03\x12\x7e\x84\xc3\xe5\xa0\xce\xca\x55\xca\x25\x5c\x8e\xc9\x39\xdd\x7d\x41\xba\x39\x17\x13\x01\x92\xa7\x22\x80\x49\x8c\x19\x0e\xeb\xe0\x77\x91\xf0\x52\x8f\xbf\x8f\x80\x5a\x28\xa7\x1e\x9a\x0b\x1e\x23\xca\x14\x08\x86\xa3\xbc\x42\x59\x46\xc8\x4b\x1a\x27\x11\x68\x1d\x1b\x10\x32\xe7\xc9\x1d\x27\x4f\x6e\x86\xb8\x9e\xac\x0a\xc8\x72\xa8\x17\x70\x9f\x82\x54\x4b\xf0\xb8\x2a\xc5\x3e\x0c\xa6\x98\x34\xa3\x30\x05\xdf\xd6\xa8\x56\x9d\xd9\xeb\xa5\x5a\x7a\x17\xd1\xa0\x83\x69\x39\xfc\xf7\x5c\xb8\xe2\xa8\x7c\xd6\x5b\xbd\x5e\x15\x82\x43\x72\x5c\xa1\xb6\xb7\x7e\xcd\xf5\x8e\x5f\xc0\x6e\x3f\xb1\xd4\x67\xbd\xd5\xfc\x52\x98\xff\xc0\xca\xbc\xe3\x3d\x65\x93\x40\x00\x56\x50\x3c\xda\x1e\xba\xff\x6b\x14\x23\x06\xdf\x0b\x90\x0f\x81\xb7\x6d\xf6\x3a\xca\xb8\x71\x12\x50\x3e\xf3\x2d\xe6\x9d\x7f\xf5\x32\xd0\xba\x9d\x20\x87\xdf\x77\x35\x89\xd7\x59\xcc\x15\xa1\xc8\x7f\xd6\xee\x7e\xd9\x7f\xa6\x79\x00\x21\xb5\x4d\x27\x7c\x5e\xc8\x5d\xb3\xed\xf7\x95\xdf\xd7\xfd\xc8\x04\x02\x3a\xa7\x41\xde\xf9\x76\x93\x5c\xf6\x22\xb7\x4a\x18\x38\xb0\x57\xb9\xb6\xd9\x76\xbb\x4a\x63\xfe\x7f\x9d\x2f\x8d\x6f\xb9\x7d\x49\x6e\x4f\xc7\x0c\x28\x69\x42\xaa\xa9\x7d\xb7\x3c\xff\x6c\x3a\x3b\xd8\x54\x3f\x32\xed\x7e\xc4\xa2\xa3\xff\x65\xd7\x02\xf2\x6f\x59\x74\x14\x21\x80\x40\x04\x0a\x7e\xd4\x28\x30\x33\xb3\x5b\x29\x0a\xd8\x26\x6f\x51\x60\x84\xb2\x63\xfb\x8c\xce\x1d\x8d\x2c\xac\x89\x77\x50\x05\x7f\x7e\xf9\xf2\x5a\x56\x73\xe1\x51\x01\x23\x40\xb6\xbf\x82\xab\x38\xc2\xec\xc9\x6e\x24\x2f\x5d\xb8\xdd\x6c\xc9\x76\x18\x6c\xbb\xe2\xff\x61\xa6\xd2\xbd\x2f\x1d\x6f\x75\xcd\xf8\x00\x57\x8b\x8b\x63\x0f\x95\x89\x94\x87\x1f\x4e\x85\xa8\xd2\xbb\x61\xc4\x53\x21\xd0\x64\x82\xa4\xc2\x8c\x60\x41\x90\x04\x41\x71\x44\xff\xc4\x77\x11\xa0\x9f\xcf\xcf\x90\x19\xef\x35\xcb\xb6\xad\xb4\x6c\xc0\x99\x16\x57\xf6\x51\xce\xb0\x93\x6b\xf6\x4f\x74\xed\x51\xf6\x80\x23\x4a\x50\x2a\x41\x68\xf8\x5d\x7b\xf6\xfe\x7d\xca\x15\x46\xf0\x18\x00\x10\x20\xf9\x5d\x23\x6b\xd8\x97\xf5\xe3\x5d\xb3\xe9\x74\x0a\x2a\x98\x4e\xa7\xd7\xec\x6c\xa6\xfb\x4b\x19\xbd\x4f\x21\xeb\x8d\x12\x60\x4a\xe7\x6e\xdb\x2a\xe0\x04\xae\xd9\x0c\x14\xa6\x91\xd4\xc2\xdc\xcc\x0c\x47\x7a\x94\x0a\x1e\x1b\x83\x94\xe8\x1b\x65\x04\xdb\xce\xe7\x14\x22\x82\x8e\x2e\xc0\x1c\x6a\x91\x47\x28\x4e\xa5\x42\x77\x80\x18\x67\x93\x3f\x41\x70\xf4\x80\xa3\xb4\x98\x01\xe3\x0a\x01\xe3\x69\xb8\x40\x8a\x86\x0b\x25\x75\x4c\x99\x03\x10\x14\xf2\x64\x01\x22\x97\xcb\x77\xa3\xd0\xd1\x47\x4e\x8e\x10\xe1\x20\x8f\x14\x82\x47\x2a\x95\x16\xf9\x45\xf7\x5a\x1f\xaa\x04\x85\xf8\x1c\x7d\x83\xa7\x89\xe9\x11\x25\x98\x8a\x72\x9d\xbd\x08\x2a\xfc\xee\x2b\x04\xaa\xbc\x9f\x08\x1d\xbc\x14\x6d\xd0\xc7\x23\xd6\x1c\x4d\x4e\xe5\x6a\xb0\x10\xb8\x4e\x7a\x8f\x2a\x88\x65\x1b\xbf\x8d\x68\xd6\x49\x80\xc7\x49\xc8\x27\x79\xb0\xc9\x9c\xe1\x39\x61\x6d\x6c\xde\x1a\x99\x9b\x8a\xd6\x52\x6e\x3d\x94\x0c\xd3\x71\x2a\xc4\xd9\xcc\xad\x22\x83\x44\x97\x95\x1c\x11\xbc\x3e\xcf\x8c\x11\x6e\xe5\x52\x61\x95\xca\x5b\x9d\x2a\xbb\x3a\xa0\x4c\x41\x7d\xb3\x52\xdb\x87\x8b\x18\xab\xec\xf1\x7f\xfe\xbd\xa4\xfb\x4b\xd3\xc3\xa7\xab\xab\xf3\xca\x08\x5a\xd5\xb6\x69\x92\xe0\xe0\x9b\x9d\xa9\x17\x52\xd5\x4a\xd7\xc1\xc2\xaf\x65\x67\xe1\x3f\x00\x23\x5c\xf8\x21\x55\x8b\xf4\x6e\x1a\xf0\xd8\xaf\xb4\xf1\x83\x05\x88\xbc\xb0\xca\xab\x0a\x6b\xe7\xa5\x61\xc6\x92\x59\x40\x22\x40\x02\x53\xb2\xce\xeb\xb3\xd9\xfa\x78\xd7\x9c\x1e\x0c\x87\x5f\x69\xb5\x74\xa8\xf9\x6c\x05\x50\x5d\xd6\x50\xb5\x4f\xb3\xff\x5a\x9f\xbc\xc3\xf0\x5a\xa2\x61\x7a\x6d\x31\x1d\x72\x6c\x32\x6a\x19\xbe\x09\xcd\x0a\x2c\xd3\x3a\x2e\xc7\x9c\xe9\x65\x1f\xc2\x2e\x5b\x10\xcb\xcf\x12\x9e\xcd\x96\xcc\xb6\xc1\xf4\xd1\xa6\xf4\x4b\x33\x16\x36\xa7\x54\xa6\x8a\xca\x94\x30\x21\x34\xcb\x1b\x65\xc2\xb0\x61\x75\xc9\x1c\x9b\x54\x2a\xb5\x9c\x77\x90\xca\x1d\xfb\xc7\x87\xb8\xe3\xbc\x4c\xb7\xc5\xda\xc2\xd6\x7a\xd9\xb5\xe2\x66\x39\x36\x2f\x1c\x45\x56\xb4\xae\x1d\x78\xb2\x4d\xa4\x5b\xac\x56\xc9\x22\x25\x99\xec\xc2\x04\x6d\xbc\x4d\xd7\xe3\xbc\x5d\x24\x24\x3f\xab\xa1\x89\x70\x70\x02\x3b\x23\x7b\x48\x8c\x39\x14\x6e\x33\xe9\xf5\xba\xc8\xd7\x4e\xfe\xa7\x7f\x3b\xfb\xd1\xf5\xe7\xed\x06\xc6\xf9\x2c\x41\x74\x19\xe8\x01\x0b\xaa\xeb\xe3\xce\x0a\xab\x81\x9f\x61\x8c\x73\x69\x18\x56\x71\x7d\x29\xc6\xb3\xf5\x94\x94\x7c\x0b\xf5\x5b\x3c\xe4\xf5\x5c\x83\x95\xbf\x51\xa9\x4e\x99\x32\xef\xdb\x3d\x9c\x2c\x44\x35\x23\x8d\x4a\x73\xc2\xbb\x76\x7c\xf9\x6f\xcb\xc4\xd9\x0f\x40\x96\xd9\x5e\xd0\xf7\x7f\xb3\xc8\x37\x08\x7d\x46\x54\xa3\xcf\x2c\x0c\x56\xd1\xb7\x21\xf8\xce\x97\x44\xcb\x41\x55\xbe\xd3\x29\xe7\x20\xb6\xa8\x78\x3c\x8f\x0c\x4f\xd3\x43\x42\x02\xfa\x4e\xd5\xc2\x7a\xec\xc7\x76\x50\x3b\x01\xdf\x1a\x54\x6e\xfc\x12\xdf\xfb\x2d\x80\x0d\xe0\xc3\x52\x4d\xfb\xfc\xfa\x48\xc8\x92\x97\x8a\x0b\x1c\xc2\x17\x5d\xef\xf6\x43\xab\x26\xee\x84\x17\xa2\x0c\x49\x2b\x74\xd0\x69\xe7\xad\x5a\x7b\xab\xd6\xfa\x38\x74\x01\x01\x4d\x28\x30\xb5\x84\x38\x85\x8c\x7d\x63\xca\x7f\x50\x36\xe7\xc5\xa7\x78\x48\xf7\x87\xec\xc2\x74\x3f\x69\x2a\x7b\x16\x7f\xd4\xca\xa4\x8a\x73\x1a\x5b\x89\x9e\xd9\x9f\xc9\x0d\x73\xd3\xcb\x3f\x2b\xbe\x2e\x6a\x4e\xab\x9d\xed\xbc\x98\xdb\x88\x42\xdd\xd4\x79\x83\xb7\x57\xdd\x80\xea\x06\x78\xfb\x6b\x80\x72\x6d\x20\x30\x6b\x03\x25\xc4\x53\x49\x59\x98\xfd\x1e\x88\xf6\xee\xc5\xf8\x08\x3f\xed\x68\x0d\x78\x66\x74\xaf\x93\x16\x5c\xa0\x70\xf2\x35\x7b\x92\x85\x83\x5b\x82\x15\xae\x02\xe1\xa6\xa6\x21\x59\x02\xa6\x80\xc7\x31\x67\xb7\x5d\x70\xed\x1b\xdd\x70\xd8\xf6\xec\x5c\xb4\xe0\xeb\xc8\xb8\x66\xa8\x0e\x20\x3b\xda\x36\x6d\xb3\x64\x56\xed\x92\x6d\x49\xd9\xb6\xa4\x74\x2b\xa3\xf9\x6a\xd3\x2a\xda\xc9\xa1\x1b\x9b\x3d\x85\xc2\x28\xb4\xee\x7d\x99\x18\x9d\xd1\x26\xe9\xec\xb4\x38\xd7\x53\xb2\xfb\x3b\x03\x83\xae\x15\x1e\x3b\xe6\x66\x43\x5c\xee\x1a\x2b\x94\xed\x31\xeb\x2b\x3e\xd7\x6e\x50\xd9\x37\xfd\xeb\xbb\xa1\x24\xde\x06\xc9\xba\xe0\x44\x57\xd6\x96\xcd\x39\xae\xa6\xbf\xe9\xc5\xd1\x8a\xeb\x5d\xbb\xbf\xf5\xd1\xe0\xb0\xc4\x6b\x64\x2d\x4b\xed\xa5\x8b\xa4\x95\x4f\x01\x37\xab\x35\x5d\x1f\xf5\xe4\xaf\x1a\xce\x12\x6a\x70\xe1\x59\x68\xde\xd4\x6d\x23\x2d\xc8\xbd\x15\x8c\x5d\xe8\xed\xcf\x2f\xab\xd7\x8d\x6b\x00\xb8\x0b\x68\x87\xb6\x34\x7b\xd5\x26\x46\xd3\xa0\xb9\x88\x36\x63\x71\x7e\xd2\xe8\x5a\x97\xe7\xcd\xb7\x48\xc7\xb1\x5e\x53\x9c\xd5\xce\x09\xd6\x30\x7a\x73\x08\x4b\x3d\x8e\x02\x73\x38\x10\x66\xb5\xca\x7d\x84\x17\xd3\xdc\xc2\x6b\xa7\xce\xb4\x11\x45\x5e\x1c\x1e\x5e\x57\xf9\x97\xac\xfd\xe8\xbc\x90\x43\xb6\xcc\xda\xc2\xf5\x25\xcc\x82\x35\x9b\xed\x5b\xec\x6a\x45\x22\xfb\x7e\x6f\x4f\xe7\xcf\xbe\xe4\xdd\xef\xc7\xb7\x03\x32\x8a\x53\xbe\x1a\x0d\xb3\xad\x8f\xd6\x99\xd6\x95\x3d\xec\xfa\xfb\x01\x5b\x7d\x53\x71\xc0\x7a\x98\x97\x1c\x9f\x54\xec\xcc\x4d\x95\x3f\x0b\x56\x39\xc5\x9d\x1d\x8f\x6d\x7c\x28\x51\xf1\x5a\x3e\xe3\xe5\x67\x56\x52\x45\x23\xe9\x2f\x94\x4a\xf4\x95\x9f\xfd\x21\xb4\xe9\x57\xc9\x59\xef\x07\x19\xe5\xc8\xda\x07\xa1\x3d\x7b\x92\xa7\x1b\x45\xf6\xb0\x4c\xf3\xc0\x8f\xe3\xbc\xf3\xc0\x89\x58\x7d\x8d\x09\x34\x8e\xbd\x39\x0e\x1e\xbf\x7b\xf9\x2b\x00\x00\xff\xff\x69\x89\xab\xa0\x59\x4e\x00\x00")

// FileVendorGithubComContainerumCherrySwaggerJSON is "vendor/github.com/containerum/cherry/swagger.json"
var FileVendorGithubComContainerumCherrySwaggerJSON = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x54\x41\x6f\xdb\x3c\x0c\xbd\xe7\x57\x10\xfe\x3e\x20\x87\x2d\xce\x30\x0c\x3b\xf4\x36\x20\xed\x16\x14\x03\x8a\xa5\x47\x03\x83\x62\xd1\x36\x57\x5b\x72\x29\xba\x4b\x5a\xf4\xbf\x0f\xb2\xdc\xc6\x76\x93\xac\x5b\x76\xf0\xc1\xd4\x23\xf9\x48\x3e\xf2\x61\x02\x10\xb9\x9f\x2a\xcf\x91\xa3\x33\x88\xde\xc7\xef\xa2\xb7\xde\x56\x2b\x29\x5c\x74\x06\x0f\x8f\xed\xaf\xc6\x8c\x0c\x09\x59\xd3\x1a\x27\x00\x00\xd1\x39\xf3\xf3\x4f\x8b\x71\x29\x53\xed\x41\x3e\xd4\x39\x33\xcc\x66\xe0\x44\x19\xad\x58\x83\x43\x26\x55\xd2\xbd\x5a\x97\x08\x9f\xae\x96\x80\xcc\x96\x13\xf3\x15\x9d\x53\x39\x7a\x6c\x6a\x8d\x87\x4b\x78\x82\x2a\xbc\x9c\x25\xe6\x0d\x24\x11\x99\x3b\x55\x92\x86\xc6\x21\x1b\x55\x61\x12\x05\xfb\x6d\x63\x45\x01\x6e\x52\x44\x8d\xfa\xc9\xda\x62\x95\xe7\xd2\xe5\x89\x12\x13\xc7\x31\x4a\x1a\xc7\x71\x62\x96\x0b\x9f\xaf\x31\x74\xdb\x60\x97\x8d\x34\x1a\xa1\x8c\xd2\xe0\x95\x5a\x8d\x89\x59\xa0\x28\x2a\x9d\x07\xdb\xb6\x32\x55\x7a\x96\x82\x9b\x11\x49\x07\x37\x64\xb4\x0a\xc9\x33\xc2\x52\xc3\xf4\x1b\xd6\x25\xa5\xca\x4d\xa1\x6a\x9c\xc0\x1a\xc1\x58\x33\xbb\x47\xb6\x70\xa7\xca\xe6\xb9\x02\x63\x05\xd0\xd8\x26\x2f\x40\x28\x2f\xc4\x81\x58\xc8\x10\x35\xe4\xb6\x2e\x90\x9f\x70\x8c\xce\x36\x9c\x22\x4c\x3f\x5b\x3d\x05\x6d\xd1\x4d\x05\x70\x43\x4e\x3c\xe4\xc2\x67\x1d\x52\x75\x28\x60\x33\xb8\xc1\xed\xac\xcd\x08\xb5\x22\x76\xed\x84\xdb\x91\xc9\xb6\x46\x3f\x2b\xbb\xfe\x81\xa9\xec\xec\x35\xdb\x1a\x59\x08\x5d\x6f\xc0\xed\x88\xdb\x76\x0c\x8c\xbd\x30\x8a\x59\x6d\x9f\xa3\xb4\x4f\x24\x58\x8d\xf1\x3d\x0f\x27\x4c\x26\x8f\x7a\x8f\x8f\x03\xf7\xcd\x2c\xb7\x33\x3f\x6d\x0f\xee\x86\xb1\x43\xf7\xb0\xa1\xe7\x2f\x98\xfd\xcf\x98\x79\xd7\xff\xe6\x3d\x05\xcf\x43\xa7\xf6\xc7\x21\xfd\xba\x18\xe7\xcc\xcb\xc5\xfe\x10\x9d\x24\x0e\x75\xa9\xab\xf9\x70\x9d\xdd\x46\xec\x0f\xee\x44\x49\xe3\xbe\x17\x22\xf5\xa1\x04\x64\x04\xfd\x3e\x0f\x32\x64\x96\x2b\x25\xdd\xf3\xc7\x0f\x47\xd2\xaf\xda\x0c\x5f\xae\xaf\xaf\x7a\x0c\x26\x23\x26\xc1\xa5\x56\xe9\x4d\xa8\x34\xca\x49\x8a\x66\x1d\xa7\xb6\x9a\xfb\x05\x51\x64\x90\x9b\x6a\x9e\x16\xc8\xbc\x0d\x81\x3a\xe7\x28\xb4\xee\xe8\xe5\x08\xfb\xc9\x58\x33\x3a\x34\xe2\x86\xab\xba\x5c\xfc\xbd\x84\xfd\x9a\xbe\x7a\xc2\x97\x1e\xbc\x7f\x0c\x7f\xa0\x93\xd5\x40\x28\xff\xb8\x93\x97\xc3\x7a\xf6\xf4\xd2\x23\x46\xdd\xf4\x4d\xf0\x87\xa1\x6d\xe7\xcb\x5e\x8e\x05\xd4\x13\x4f\x33\x54\xcf\x89\xe4\x57\xbf\xd3\xc1\xea\x85\x10\x1c\xf2\x1d\xa5\x08\xcb\xc5\x91\x02\x46\x2b\x76\x0a\xcb\x8b\xf1\x5d\x19\xb3\xdc\x9d\xdd\x1e\x4b\xa5\x35\x75\x37\x78\x77\x7c\xc3\x89\x3a\x42\x7b\xac\xe1\x5d\x94\xab\x03\x6a\xde\x7f\x47\x4f\xd3\xd6\xc4\x7f\x8f\xbf\x02\x00\x00\xff\xff\x85\xa0\x34\xb4\x1a\x08\x00\x00")

// FileVendorGithubComContainerumUtilsHttputilSwaggerJSON is "vendor/github.com/containerum/utils/httputil/swagger.json"
var FileVendorGithubComContainerumUtilsHttputilSwaggerJSON = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xbc\x95\xcf\x6a\xe3\x30\x10\xc6\xef\x79\x8a\x41\xe7\xf5\xb2\xb0\x37\xdf\x96\xf5\x21\xbe\x2c\x21\x9b\x96\x42\xe9\x61\x62\x8f\x6d\x51\x5b\x72\xf4\x87\x10\x4a\xde\xbd\x48\x72\x9d\x86\x94\xc8\x4d\x71\x0f\x21\x68\xe6\xf3\xcc\x4f\x9f\xc6\xf2\xcb\x02\x80\xe9\x3d\xd6\x35\x29\x96\x02\xfb\xfd\xf3\x17\xfb\xe1\x62\x3d\x2a\xec\xc8\x90\xd2\x2c\x05\xa7\x02\x60\xff\xed\x56\x1b\x6e\xac\xa1\xf2\x4e\x93\xca\xb3\x31\x05\xc0\x04\x76\xe4\x2a\x58\x4d\x2a\xe1\xa5\xaf\xe2\x13\x5c\xb8\xf0\xce\x92\x3a\x9c\x82\x8a\x76\x96\x2b\x2a\x59\x0a\x15\xb6\x9a\xc6\x84\x39\xf4\xbe\x8c\x36\x8a\x8b\xfa\xf4\x40\x25\x55\x87\xc6\x37\xb0\xef\xab\x97\xa4\x0b\xc5\x7b\xc3\xa5\x6f\xb3\x69\x08\xf2\x0c\x64\x05\xa6\xa1\xb7\x3f\x87\x04\x52\xc0\x96\x1a\x6c\x2b\x17\xdd\x37\xbc\x68\x7c\x0e\xcb\x8e\x0b\xae\x8d\x42\x23\x15\x14\x28\xa0\x27\xe5\x9a\x85\x6c\xe1\x0b\xfb\x66\xc7\xd0\x93\xad\x69\x67\x49\x9b\x3c\x5b\x12\x96\xde\xb5\x0b\x0f\x1e\x92\x41\x94\xe4\x19\xbb\x79\x6b\xc1\xb8\x26\x74\xf9\xc8\x39\xa3\x2c\x9d\xa1\x85\x53\xb9\xc6\xe5\x14\xdf\x0c\x35\x58\xf1\x0f\x3b\x9a\xe2\x98\xd3\xc5\xf1\x6e\xf3\xe6\x6f\xcb\x49\x98\xa8\x3f\x41\x36\x17\xc4\x9f\x7a\x0a\x83\x57\xcd\x81\x10\x36\x97\xaf\xae\x11\x04\x4d\x92\xaf\x6e\x05\x38\x4d\x0f\xef\xa7\x61\x6d\xe4\x33\x89\xeb\xb3\xeb\x25\x93\x86\x37\x06\x75\x3e\xd2\xb1\x03\x5b\xcb\x96\xa2\xe7\xe5\x44\x71\x30\x12\xb6\x63\x29\x3c\x0e\x6b\x08\xd7\xe5\x98\x07\x60\xfe\x3e\x62\xc3\xfa\xe9\x8b\xa3\xe6\xde\x25\xdd\x63\x41\x3a\xca\x3f\x4a\xe7\x1a\xfb\x7b\xd9\xda\x6e\x02\x48\xd0\xcd\x45\xb1\xe4\x25\x65\x68\x30\x8a\xe1\x84\x89\x53\x5e\x92\x6c\xa5\x6c\x09\xc5\x27\x50\xfc\x27\x2e\xb0\x2c\xdc\xef\xf8\x1a\x00\x00\xff\xff\xce\xd4\x03\x4a\x77\x07\x00\x00")

func init() {
	if CTX.Err() != nil {
		panic(CTX.Err())
	}

	var err error

	err = FS.Mkdir(CTX, "vendor/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	err = FS.Mkdir(CTX, "vendor/github.com/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	err = FS.Mkdir(CTX, "vendor/github.com/containerum/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	err = FS.Mkdir(CTX, "vendor/github.com/containerum/utils/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	err = FS.Mkdir(CTX, "vendor/github.com/containerum/utils/httputil/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	err = FS.Mkdir(CTX, "vendor/github.com/containerum/cherry/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	var f webdav.File

	var rb *bytes.Reader
	var r *gzip.Reader

	rb = bytes.NewReader(FileSwaggerJSON)
	r, err = gzip.NewReader(rb)
	if err != nil {
		panic(err)
	}

	err = r.Close()
	if err != nil {
		panic(err)
	}

	f, err = FS.OpenFile(CTX, "swagger.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	rb = bytes.NewReader(FileVendorGithubComContainerumCherrySwaggerJSON)
	r, err = gzip.NewReader(rb)
	if err != nil {
		panic(err)
	}

	err = r.Close()
	if err != nil {
		panic(err)
	}

	f, err = FS.OpenFile(CTX, "vendor/github.com/containerum/cherry/swagger.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	rb = bytes.NewReader(FileVendorGithubComContainerumUtilsHttputilSwaggerJSON)
	r, err = gzip.NewReader(rb)
	if err != nil {
		panic(err)
	}

	err = r.Close()
	if err != nil {
		panic(err)
	}

	f, err = FS.OpenFile(CTX, "vendor/github.com/containerum/utils/httputil/swagger.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	Handler = &webdav.Handler{
		FileSystem: FS,
		LockSystem: webdav.NewMemLS(),
	}

}

// Open a file
func (hfs *HTTPFS) Open(path string) (http.File, error) {

	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(f)
	return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := FS.OpenFile(CTX, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// WalkDirs looks for files in the given dir and returns a list of files in it
// usage for all files in the b0x: WalkDirs("", false)
func WalkDirs(name string, includeDirsInList bool, files ...string) ([]string, error) {
	f, err := FS.OpenFile(CTX, name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	fileInfos, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	for _, info := range fileInfos {
		filename := path.Join(name, info.Name())

		if includeDirsInList || !info.IsDir() {
			files = append(files, filename)
		}

		if info.IsDir() {
			files, err = WalkDirs(filename, includeDirsInList, files...)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}
