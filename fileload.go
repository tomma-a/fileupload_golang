package main
import (
	"fmt"
	"net/http"
	"flag"
	"io"
	"os"
	"path"
	 "embed"
)
//go:embed t.html
var htext string

//go:embed bootstrap/*
var efs embed.FS
var updir string="/tmp/upload/"
func upload_handle(w http.ResponseWriter, r *http.Request) {
	if r.Method=="GET" {
		fmt.Fprintf(w,htext)
	} else {
		r.ParseMultipartForm(32 <<20)
		files :=r.MultipartForm.File["uploadfile"]
		for i,_ :=range files {
		ff ,err:=files[i].Open()
		if err!=nil {
			panic(err)
			}
		defer ff.Close()
		f, err:=os.OpenFile(path.Join(updir,files[i].Filename),os.O_WRONLY|os.O_CREATE,0666)
		if err !=nil {
			panic(err)
		}
		defer f.Close()
		io.Copy(f,ff);
		fmt.Fprintf(w,"Successfully uploaded file%s!\n",files[i].Filename);
		}
	}
}
func ensuredir(s string) {
	finfo,err :=os.Stat(s)
	if err!=nil {
		if os.IsNotExist(err)	 {
			os.Mkdir(s,0750)
		} else {
			panic(err)
		}
	}
	if !finfo.IsDir()  {
		panic(s+" is not a dir")
	}
}
func main() {
	port:=flag.String("p","8100","port to listen")
	dir:=flag.String("d","/tmp","directory to save upload file")
	flag.Parse()
	updir=*dir
	ensuredir(updir)
	fs:=http.FileServer(http.FS(efs))
	http.Handle("/static/",http.StripPrefix("/static/",fs))
	http.HandleFunc("/upload",upload_handle)
	fmt.Printf("Preparing to run :%s .....\n",*port)
	http.ListenAndServe(":"+*port,nil)
}
