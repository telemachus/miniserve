# miniserve: a small web server for local development

## What?

`miniserve` is a small HTTP server for testing websites on a development
computer.  It doesn't do much, but it has a few features I wanted.

+ URLs don't require `.html`.  For example, `miniserve` will look for
  `basedir/pretty-url` and, failing that, will try to serve
  `basedir/pretty-url.html`.
+ Structured logging by default.
+ Graceful shutdown, as much as possible.
+ Optional flags can specify what port to serve on and what directory to
  serve.  (The defaults are 8080 and `"."`.)

## Why?

[My website](https://telemachus.me) uses pretty URLs (i.e., no `.html` at the
end).  On the server, `nginx` takes care of all that, but I needed a way to
test the site at home.  Once upon a time, I wrote a very ugly solution, but it
had no logging.  I added logging using
[Gorilla](https://github.com/gorilla/handlers), but Gorilla recently [stopped
being maintained](https://github.com/gorilla#gorilla-toolkit).  I took the
opportunity to write a fuller and better—though still miniature—server.

To be honest, I wrote this largely as an excuse to learn things. It does what
I need well, but I wouldn't recommend it for anything serious. I wouldn't even
recommend it to anyone but me even for trivial things. That said, if you have
suggestions or criticisms, please file an issue.

## Credits

Thanks to the authors of the following for inspiration and (often) code.

+ [Golang HTTP Server Graceful Shutdown](https://clavinjune.dev/en/blogs/golang-http-server-graceful-shutdown)
+ [http-server-shutdown.go](https://gist.github.com/rgl/0351b6d9362abb32d6b55f86bd17ab65)
+ [Serve pages from links without `.html`](https://stackoverflow.com/a/57281956)
+ [A Guide to Writing Logging Middleware in Go](https://blog.questionable.services/article/guide-logging-middleware-go)
+ [How I write HTTP services](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html)
+ [The http.Handler wrapper technique in #golang UPDATED](https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702)
+ [Understanding http.HandlerFunc](https://stackoverflow.com/q/53678633)
+ [logfmt](https://brandur.org/logfmt)
+ [Using Canonical Log Lines for Online Visibility](https://brandur.org/canonical-log-lines)
