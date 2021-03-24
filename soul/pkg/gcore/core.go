package gcore

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"google.golang.org/grpc/credentials"
)

//go-grpc-middleware 拦截器中间件 在网关加各种判断

//微服务配置项
type ServeSetting struct {
	Host   string
	Server     interface{}
	Register   interface{}
	HandlerFromEndpoint   interface{}
}


//启动http服务
func RunServeHTTP(servers []ServeSetting,port string) {
	PORT := ":"+port
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	//把server注册HTTP
	for _,inter := range servers{
		fn:= reflect.ValueOf(inter.HandlerFromEndpoint)
		params := make([]reflect.Value, 4)
		params[0] = reflect.ValueOf(ctx)
		params[1] = reflect.ValueOf(mux)
		params[2] = reflect.ValueOf(inter.Host)
		params[3] = reflect.ValueOf(opts)
		fn.Call(params)
	}

	log.Printf("listen http on "+PORT)

	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return

}

//启动grpc服务
func RunServeGRPC(servers []ServeSetting,port string) {
	PORT := ":"+port
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("listen grpc on "+PORT)

	server := grpc.NewServer()

	for _,inter := range servers{
		fn:= reflect.ValueOf(inter.Register)
		params := make([]reflect.Value, 2)
		params[0] = reflect.ValueOf(server)
		params[1] = reflect.ValueOf(inter.Server)
		fn.Call(params)
	}

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}





var (
	demoKeyPair  *tls.Certificate
	demoCertPool *x509.CertPool
)

func MakeInsecure (key string,cert string) {
	var err error
	pair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		panic(err)
	}
	demoKeyPair = &pair
	demoCertPool = x509.NewCertPool()
	ok := demoCertPool.AppendCertsFromPEM([]byte(cert))
	if !ok {
		panic("bad certs")
	}
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}


//	ServerMap["test"] = map[string]interface{}{"Host":"localhost:10000","Server":&service.TestServer{},"Register":pb.RegisterTestServer,"HandlerFromEndpoint":pb.RegisterTestHandlerFromEndpoint}
func RunServe(servers []ServeSetting,serveAddr string,port int) {

	//启动GRPC服务
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(demoCertPool, serveAddr))}

	grpcServer := grpc.NewServer(opts...)

	//把server注册GRPC
	for _,inter := range servers{
		fn:= reflect.ValueOf(inter.Register)
		params := make([]reflect.Value, 2)
		params[0] = reflect.ValueOf(grpcServer)
		params[1] = reflect.ValueOf(inter.Server)
		fn.Call(params)
	}

	//启动网关，把grpc注册HTTP
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: serveAddr,
		RootCAs:    demoCertPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		//io.Copy(w, strings.NewReader(pb.Swagger))
	})

	gwmux := runtime.NewServeMux()

	//把server注册HTTP
	for _,inter := range servers{
		fn:= reflect.ValueOf(inter.HandlerFromEndpoint)
		params := make([]reflect.Value, 4)
		params[0] = reflect.ValueOf(ctx)
		params[1] = reflect.ValueOf(gwmux)
		params[2] = reflect.ValueOf(inter.Host)
		params[3] = reflect.ValueOf(dopts)
		fn.Call(params)
	}

	mux.Handle("/", gwmux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    serveAddr,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*demoKeyPair},
			NextProtos:   []string{"h2"},
		},
	}

	fmt.Printf("grpc on port: %d\n", port)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}