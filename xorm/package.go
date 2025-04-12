package xorm

//go:generate mockgen -package xorm -destination=./engine_mock.go xorm.io/xorm EngineInterface
//go:generate mockgen -package xorm -destination=./session_mock.go -source transaction.go XInterface
