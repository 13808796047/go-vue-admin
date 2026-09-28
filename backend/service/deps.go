package service

// Deps 路由器依赖注入
type Deps struct {
	Svc       *Service
	JWTSecret string
}
