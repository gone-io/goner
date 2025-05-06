package gin

//go:generate mockgen -package=gin -destination=net_mock.go net Listener,Conn,Addr

//go:generate mockgen -package=gin -destination=http_mock.go net/http Handler

//go:generate mockgen -package=gin  -destination=origin_mock.go github.com/gin-gonic/gin ResponseWriter

//go:generate mockgen -package=gin  -destination=gin_mock.go github.com/gone-io/goner/gin Responser,HttInjector,XContext,HandleProxyToGin,Middleware,Controller
