package followredirect

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"net/http"
)

func init() {
	caddy.RegisterModule(FollowRedirect{})
	httpcaddyfile.RegisterHandlerDirective("follow_redirect", parseCaddyfile)
}

// FollowRedirect 定义结构体,也定义了所需参数之类的
type FollowRedirect struct {
	Recursive bool `json:"recursive,omitempty"`
	logger    *zap.Logger
	Transport http.RoundTripper `json:"-"`
}

// Provision 就是解析然后往结构体塞数据
func (f FollowRedirect) Provision(c caddy.Context) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

// ServeHTTP 在这里处理 HTTP 内容
func (f FollowRedirect) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	//repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	logger := f.logger.With(
		zap.Object("request", caddyhttp.LoggableHTTPRequest{Request: r}),
	)
	logger.Debug(w.Header().Get("Location"))

	//fmt.Println()
	return next.ServeHTTP(w, r)
}

func (f FollowRedirect) CaddyModule() caddy.ModuleInfo {
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
	var m FollowRedirect
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
