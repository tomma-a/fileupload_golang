# This is just an example to get you started. A typical binary package
# uses this file as the main entry point of the application.

import jester,os,strutils,asyncdispatch,zippy,httpcore,htmlgen
const htext=staticRead "../t.html"
const bootstrap_css=staticRead "../bootstrap/bootstrap.min.css"
const bootstrap_js=staticRead "../bootstrap/bootstrap.min.js"
router myroute:
        get "/static/bootstrap/bootstrap.min.css":
                if request.headers.hasKey("accept-encoding") and ($request.headers["accept-encoding"]).contains("gzip"):
                        resp(Http200,[("Content-Encoding", "gzip")],compress(bootstrap_css,BestCompression,dfGzip))
                else:
                        resp bootstrap_css
        get "/static/bootstrap/bootstrap.min.js":
                if request.headers.hasKey("accept-encoding") and ($request.headers["accept-encoding"]).contains("gzip"):
                        resp(Http200,[("Content-Encoding", "gzip")],compress(bootstrap_js,BestSpeed,dfGzip))
                else:
                        resp bootstrap_js
        get "/upload":
                resp htext
        post "/upload":
                var msg=""
                for k,v in request.formData:
                        if k=="uploadfile":
                                writeFile("/tmp/" & v.fields["filename"],v.body)
                                msg=msg &  "write file " & v.fields["filename"] & "successfully\n"
                resp "ok\n" & msg
        get "/":
                var html=""
                for ff in walkDir("."):
                        html.add(li(a(href=ff.path[2..^1],ff.path[2..^1])))
                        html.add(br())
                resp html
proc main() = 
        let port=paramStr(1).parseInt().Port
        let settings=newSettings(port=port,staticDir=getCurrentDir())
        var jester=initJester(myroute,settings)
        jester.serve()
when isMainModule:
        main()
