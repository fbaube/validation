package validation

import(
	"embed"
	"io"
	"io/fs"
	"os"
	"fmt"
	S "strings"
	FP "path/filepath"
	)

//go:embed oasis-lwdita
var OasisLwditaFS embed.FS
// data, _ := f.ReadFile("hello.txt")

var Sep = string(os.PathSeparator)

func init() {
     ListFS("oasis-lwdita",&OasisLwditaFS)
     cache = make(map[string]string)
     }

// cache maps simple file names to redd-in file contents.
var cache map[string]string

/*
https://pkg.go.dev/embed#hdr-File_Systems

FS implements interface io/fs.FS, so it can be used
with any package that understands file systems,
including net/http, text/template, html/template.

func (f FS) Open(name string) (fs.File, error)
   (The returned file implements io.Seeker and
     io.ReaderAt when the file is not a directory)
func (f FS) ReadDir(name string) ([]fs.DirEntry, error)
func (f FS) ReadFile(name string) ([]byte, error)

For example,
http.Handle(
	"/static/",
	http.StripPrefix("/static/",
	http.FileServer(http.FS(content))))
Or,
	content1, _ := folder.ReadFile("folder/file1.hash")
	print(string(content1))
*/

func ListFS(name string, pFS *embed.FS) error {
     var DE []fs.DirEntry
     var e error
     
     DE, e = pFS.ReadDir(".")
     if e != nil {
     	return fmt.Errorf("ListFS(.): %w", e)
	}
     fmt.Printf("m5app/%s: %+v \n", name, DE)
/*
     DE, e = pFS.ReadDir("embedfs/h-body")
     if e != nil {
     	return fmt.Errorf("ListFS(h-body): %w", e)
	}
     fmt.Printf("h-body: \n%+v \n", DE)

     DE, e = pFS.ReadDir("embedfs/h-frag")
     if e != nil {
     	return fmt.Errorf("ListFS(h-frag): %w", e)
	}
     fmt.Printf("h-frag: \n%+v \n", DE)
*/
     return nil
     }

func GetContents(path string) (string, error) {
     var s string
     var e error 
     var pathSubdir string
     var pathFilext string
     var newPath string 
     var ok bool
     var f fs.File 
     
     // Try the path as-is
     if s,ok = cache[path]; ok {
     	return s, nil
	}
	
     // Analyse the path
     // Does it include a file extension ? 
     if fext := FP.Ext(path); fext != "" && len(fext) >= 2 {
     	// Form a subdir name from it (omitting the dot) 
	pathFilext = S.ToLower(fext[1:])
	println("path filext:", pathFilext)
	}
     // Does it include a subdirectory ? 
     if S.Contains(path, Sep) {
	pathSubdir = S.ToLower(path[:S.Index(path, Sep)])
	println("path subdir:", pathSubdir)
	}
	
     // If it has a file extension but no subdir, try prepending
     // the file extension (lowercased) to the path as a subdir 
     if pathFilext != "" && pathSubdir == "" {
	newPath = S.ToLower(pathFilext) + Sep + path
	if s,ok = cache[newPath]; ok {
	   return s, nil
	   }
	}
     // If it has a subdir but no file extension, try appending
     // the subdir (lowercased) to the path as a file extension 
     if pathSubdir != "" && pathFilext == "" { 
	newPath = path + "." + S.ToLower(pathFilext)
	if s,ok = cache[newPath]; ok {
	   return s, nil
	   }
	}
	
     // Here we know that we do not have it cached, so fetch
     // it out of [EmbedFS]. Use newPath, which should have 
     // all the filepath parts we use for naming & organising. 
     f,e = OasisLwditaFS.Open("oasis-lwdita" + Sep + newPath) 
     if e != nil {
     	if newPath != "" {
	   newPath = " (or " + newPath + ")"
	   return "", fmt.Errorf(
	   	"rest.embed.GetContents(%s%s): %w", path, newPath, e)
 	   }
	}
     bb, ee := io.ReadAll(f)
     s = string(bb)
     cache[path] = s
     if newPath != "" { cache[newPath] = s } 
     return s, ee
     }

