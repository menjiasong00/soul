package gcore

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
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

	//http扩展
	var gwOpts = []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),//修复json解析默认值问题
		//...
	}

	PORT := ":"+port
	ctx := context.Background()
	mux := runtime.NewServeMux(gwOpts...)
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

	if err := http.ListenAndServe(PORT, AllowCORS(mux)); err != nil {
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


	//GRPC 扩展
	var servOpts = []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second,   // If a client pings more than once every 5 seconds, terminate the connection
				PermitWithoutStream: true,              // Allow pings even when there are no active streams
			}),
		grpc.MaxRecvMsgSize(1024 * 1024 * 20),    // 新建服务器，注册服务，防恐,调整, 修改grpc默认接收的msg大小
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,  // If a client is idle for 15 seconds, send a GOAWAY
			MaxConnectionAge:      300 * time.Second, // If any connection is alive for more than 300 seconds, send a GOAWAY
			MaxConnectionAgeGrace: 5 * time.Second,   // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
			Time:                  5 * time.Second,   // Ping the client if it is idle for 5 seconds to ensure the connection is still active
			Timeout:               1 * time.Second,   // Wait 1 second for the ping ack before assuming the connection is dead
		}),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			//grpc_prometheus.StreamServerInterceptor,
			//grpc_zap.StreamServerInterceptor(zapLogger),
			//grpc_auth.StreamServerInterceptor(myAuthFunction),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			//grpc_prometheus.UnaryServerInterceptor,
			//grpc_zap.UnaryServerInterceptor(zapLogger),
			//grpc_auth.UnaryServerInterceptor(myAuthFunction),
			grpc.UnaryServerInterceptor(checkauth),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}


	server := grpc.NewServer(servOpts...)

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

func checkauth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok != true {
		return "", grpc.Errorf(codes.Unauthenticated, "found")
	}
	if len(md["grpcgateway-authorization"]) == 0 {

		return "",   grpc.Errorf(codes.Unauthenticated, "authorization no found")
	}

	split := strings.Split(md["grpcgateway-authorization"][0], " ")
	// 解码认证token
	decodeBytes, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return "", err
	}

	user := strings.Split(string(decodeBytes), ":")
	username := user[0]
	password := user[1]
	fmt.Println(username,password)
/*


	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil,grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return nil,grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}
*/
    // 继续处理请求
    return handler(ctx, req)
}







// Don't do this without consideration in production systems.
func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	//headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", "*")
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
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