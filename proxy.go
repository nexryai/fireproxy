package fireproxy

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func reverseProxy(targetBase string) gin.HandlerFunc {
	targetURL, err := url.Parse(targetBase)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		// プロキシ用のディレクター関数を設定
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = strings.TrimSuffix(targetURL.Path, "/") + c.Param("authPath")
			req.Host = targetURL.Host
		}

		// プロキシエラーの処理
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			http.Error(w, "Proxy error: "+e.Error(), http.StatusBadGateway)
		}

		// プロキシを実行
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func ConfigFirebaseAuthenticationProxy(router *gin.Engine, project string) {
	firebaseProjectDomain := project + ".firebaseapp.com"
	router.Any("/__/auth/*authPath", reverseProxy("https://"+firebaseProjectDomain+"/__/auth"))
}
