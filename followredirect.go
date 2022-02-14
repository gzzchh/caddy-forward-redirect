package followredirect

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func init() {
	caddy.RegisterModule(FollowRedirect{})
	httpcaddyfile.RegisterHandlerDirective("follow_redirect", parseCaddyfile)
}

// FollowRedirect 定义结构体,也定义了所需参数之类的
type FollowRedirect struct {
	Recursive           bool `json:"recursive,omitempty"`
	logger              *zap.Logger
	ctx                 caddy.Context
	TransportRaw        http.RoundTripper `json:"-"`
	ReverseProxyHandler reverseproxy.Handler
}

// Provision 就是解析然后往结构体塞数据
func (f *FollowRedirect) Provision(ctx caddy.Context) error {
	// start by loading modules
	if f.TransportRaw != nil {
		mod, err := ctx.LoadModule(f, "TransportRaw")
		if err != nil {
			return fmt.Errorf("loading transport: %v", err)
		}
		f.TransportRaw = mod.(http.RoundTripper)
	}
	f.ReverseProxyHandler = reverseproxy.Handler{}
	err := f.ReverseProxyHandler.Provision(ctx)
	if err != nil {
		f.logger.Error("ReverseProxyHandler Provision error", zap.Error(err))
		return err
	}
	f.ctx = ctx
	f.logger = ctx.Logger(f)
	return nil
}

// ServeHTTP 在这里处理 HTTP 内容
func (f *FollowRedirect) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	newLocation, _ := repl.Get("http.reverse_proxy.header.Location")
	newLocationStr := newLocation.(string)
	//fmt.Println(newLocation)
	newUrl, err := url.ParseRequestURI(newLocationStr)
	newUrlPort := 0
	if strings.HasPrefix(newLocationStr, "https") {
		newUrlPort = 443
	} else {
		newUrlPort = 80
	}
	if err != nil {
		f.logger.Error("url parse error", zap.Error(err))
	}
	// 在这里记录原始的请求
	origURLScheme := r.URL.Scheme
	origURLHost := r.URL.Host
	// 别忘了恢复哦
	defer func() {
		r.URL.Scheme = origURLScheme
		r.URL.Host = origURLHost
	}()

	// 改写 URL
	r.URL = newUrl
	r.RequestURI = newUrl.RequestURI()
	r.Host = newUrl.Host + ":" + strconv.Itoa(newUrlPort)
	r.TLS.ServerName = newUrl.Host
	//fmt.Println(newUrl.Scheme, newUrl.Host)

	// 下面的初始化是构造最简化能直接调用反代的代码
	// 进行一个 Upstreams 的初始化
	f.ReverseProxyHandler.Upstreams = make(reverseproxy.UpstreamPool, 1)

	// 进行一个 LoadBalancing 的初始化
	f.ReverseProxyHandler.LoadBalancing = &reverseproxy.LoadBalancing{
		SelectionPolicy: &reverseproxy.RandomSelection{},
	}
	// 进行一个 Upstream 的初始化
	f.ReverseProxyHandler.Upstreams[0] = &reverseproxy.Upstream{
		Host: &upstreamHost{},
		Dial: newUrl.Host + ":" + strconv.Itoa(newUrlPort),
	}
	// 进行一个 TransportRaw 的初始化
	if f.ReverseProxyHandler.Transport == nil {
		t := &reverseproxy.HTTPTransport{
			KeepAlive: &reverseproxy.KeepAlive{
				ProbeInterval:       caddy.Duration(30 * time.Second),
				IdleConnTimeout:     caddy.Duration(2 * time.Minute),
				MaxIdleConnsPerHost: 32, // seems about optimal, see #2805
			},
			DialTimeout: caddy.Duration(10 * time.Second),
			TLS:         &reverseproxy.TLSConfig{},
			Versions:    []string{"1.1", "2"},
			Transport:   &http.Transport{},
		}
		f.ReverseProxyHandler.Transport = t
	}
	//reflect.ValueOf(&f.ReverseProxyHandler).FieldByName("logger").SetPointer(unsafe.Pointer(&f.logger))

	// 进行一个 HealthChecks 的初始化
	//if f.ReverseProxyHandler.HealthChecks == nil {
	//	f.ReverseProxyHandler.HealthChecks = new(reverseproxy.HealthChecks)
	//}
	//if f.ReverseProxyHandler.HealthChecks.Active == nil {
	//	f.ReverseProxyHandler.HealthChecks.Active = new(reverseproxy.ActiveHealthChecks)
	//}

	f.ReverseProxyHandler.ServeHTTP(w, r, next)

	//logger.Debug(w.Header().Get("Location"))
	//fmt.Println()
	return next.ServeHTTP(w, r)
}

func (FollowRedirect) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.follow_redirect",
		New: func() caddy.Module { return new(FollowRedirect) },
	}
}

// UnmarshalCaddyfile 解析指令
func (f *FollowRedirect) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	m := new(FollowRedirect)
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner = (*FollowRedirect)(nil)
	//_ caddy.Validator             = (*FollowRedirect)(nil)
	_ caddyhttp.MiddlewareHandler = (*FollowRedirect)(nil)
	_ caddyfile.Unmarshaler       = (*FollowRedirect)(nil)
)
