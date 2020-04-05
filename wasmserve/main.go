// wasmserve compiles a Go command to WebAssembly and serves it via HTTP.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/shurcooL/go/osutil"
	"github.com/shurcooL/home/httputil"
	"github.com/shurcooL/httpgzip"
	"github.com/shurcooL/users"
)

var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address.")

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	err := http.ListenAndServe(*httpFlag, httputil.ErrorHandler(admin{}, handler{}.ServeHTTP))
	return err
}

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
	if err := httputil.AllowMethods(req, http.MethodGet, http.MethodHead); err != nil {
		return err
	}
	switch {
	case req.URL.Path == "/main.wasm":
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
		wasmFile := filepath.Join(tempDir, "main.wasm")
		cmd := exec.CommandContext(req.Context(), "go", "build", "-tags=nethttpomithttp2", "-o", wasmFile)
		env := osutil.Environ(os.Environ())
		env.Set("GOOS", "js")
		env.Set("GOARCH", "wasm")
		cmd.Env = env
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%q: %v\n\n%s", cmd.Args, err, out)
		}
		w.Header().Set("Content-Type", "application/wasm")
		return serveFile(w, req, wasmFile)
	case req.URL.Path == "/wasm_exec.js":
		w.Header().Set("Content-Type", "application/javascript")
		return serveFile(w, req, filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
	case req.URL.Path == "/favicon.ico":
		return os.ErrNotExist
	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if req.Method == http.MethodHead {
			return nil
		}
		_, err := io.WriteString(w, `<!DOCTYPE html>
<html lang="en">
	<head>
		<script src="/wasm_exec.js"></script>
		<script>
			if (!WebAssembly.instantiateStreaming) { // polyfill for Safari :/
				WebAssembly.instantiateStreaming = async (resp, importObject) => {
					const source = await (await resp).arrayBuffer();
					return await WebAssembly.instantiate(source, importObject);
				};
			}
			const go = new Go();
			const resp = fetch("/main.wasm").then((resp) => {
				if (!resp.ok) {
					resp.text().then((body) => {
						document.body.innerHTML = "<pre>" + body + "</pre>";
					});
					throw new Error("did not get acceptable status code: " + resp.status);
				}
				return resp;
			});
			WebAssembly.instantiateStreaming(resp, go.importObject).then((result) => {
				go.run(result.instance);
			}).catch((error) => { document.body.textContent = error; });
			window.addEventListener('keydown', (event) => {
				if (event.key !== '®') {
					return;
				}
				event.preventDefault();
				document.body.innerHTML = "";
				const resp = fetch("/main.wasm").then((resp) => {
					if (!resp.ok) {
						resp.text().then((body) => {
							document.body.innerHTML = "<pre>" + body + "</pre>";
						});
						throw new Error("did not get acceptable status code: " + resp.status);
					}
					return resp;
				});
				WebAssembly.instantiateStreaming(resp, go.importObject).then((result) => {
					go.run(result.instance);
				}).catch((error) => { document.body.textContent = error; });
			});
		</script>
	</head>
	<body></body>
</html>`)
		return err
	}
}

func serveFile(w http.ResponseWriter, req *http.Request, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	httpgzip.ServeContent(w, req, fi.Name(), fi.ModTime(), f)
	return nil
}

type admin struct{}

func (admin) GetAuthenticated(context.Context) (users.User, error) {
	return users.User{SiteAdmin: true}, nil
}
