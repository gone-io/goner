package gin

//go:generate mockgen -package=gin -destination=mock_net_test.go net Listener,Conn,Addr

//go:generate mockgen -package=gin -destination=mock_http_test.go net/http Handler

//go:generate mockgen -package=gin  -destination=mock_gone_test.go github.com/gone-io/gone/v2 Logger,Loader,FuncInjector

//go:generate mockgen -package=gin  -destination=mock_origin_test.go github.com/gin-gonic/gin ResponseWriter

//go:generate mockgen -package=gin  -destination=mock_g_test.go github.com/gone-io/goner/g Cmux,Tracer

//go:generate mockgen -package=gin  -destination=mock_gin_test.go -self_package=github.com/gone-io/goner/gin -source=interface.go
