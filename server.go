package main

import (
	"bytes"
	gzip "compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"sync"

	klaGzip "github.com/klauspost/compress/gzip"
)

var (
	originBuff []byte
	spWriter   sync.Pool
	spBuffer   sync.Pool
)

func main() {
	go InitPprof()

	http.HandleFunc("/oldGzip", OldGzip)
	http.HandleFunc("/klaGzip", KlaGzip)
	http.HandleFunc("/myGzip", MyGzip)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func init() {
	var err error
	originBuff, err = ioutil.ReadFile("./t.txt")
	if err != nil {
		panic(err)
	}

	// 公共对象池,更极致的优化可以建多个池
	spWriter = sync.Pool{New: func() interface{} {
		buf := new(bytes.Buffer)
		return gzip.NewWriter(buf)
	}}
	spBuffer = sync.Pool{New: func() interface{} {
		return new(bytes.Buffer)
	}}
}

func InitPprof() {
	pprofServeMux := http.NewServeMux()
	pprofServeMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofServeMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pprofServeMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofServeMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	if err := http.ListenAndServe("0.0.0.0:8081", pprofServeMux); err != nil {
		panic(err)
	}
}

func OldGzip(wr http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	w := gzip.NewWriter(buf)

	leng, err := w.Write(originBuff)
	if err != nil || leng == 0 {
		return
	}
	err = w.Flush()
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	b := buf.Bytes()
	wr.Write(b)

	// 查看是否兼容go官方gzip
	/*gr, _ := gzip.NewReader(buf)
	defer gr.Close()
	rBuf, err := ioutil.ReadAll(gr)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rBuf))*/
}

func KlaGzip(wr http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	w := klaGzip.NewWriter(buf)

	leng, err := w.Write(originBuff)
	if err != nil || leng == 0 {
		return
	}
	err = w.Flush()
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	b := buf.Bytes()
	wr.Write(b)

	// 查看是否兼容go官方gzip
	/*gr, _ := gzip.NewReader(buf)
	defer gr.Close()
	rBuf, err := ioutil.ReadAll(gr)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rBuf))*/

}

func MyGzip(wr http.ResponseWriter, r *http.Request) {
	buf := spBuffer.Get().(*bytes.Buffer)
	w := spWriter.Get().(*gzip.Writer)
	w.Reset(buf)
	defer func() {
		// 归还buff
		buf.Reset()
		spBuffer.Put(buf)
		// 归还Writer
		spWriter.Put(w)
	}()

	leng, err := w.Write(originBuff)
	if err != nil || leng == 0 {
		return
	}
	err = w.Flush()
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	b := buf.Bytes()
	wr.Write(b)

	// 查看是否兼容go官方gzip
	/*gr, _ := gzip.NewReader(buf)
	defer gr.Close()
	rBuf, err := ioutil.ReadAll(gr)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rBuf))*/
}
