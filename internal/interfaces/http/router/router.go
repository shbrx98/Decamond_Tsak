package router

import (
    "net/http"

    "github.com/gin-gonic/gin"

    appAuth "github.com/shbrx98/Decamond_Tsak/internal/application/auth"
    appUser "github.com/shbrx98/Decamond_Tsak/internal/application/user"
    "github.com/shbrx98/Decamond_Tsak/internal/infrastructure/config"
    "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/handlers"
    "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/middleware"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
)

func NewRouter(
    requestOTP *appAuth.RequestOTPUseCase,
    verifyOTP *appAuth.VerifyOTPUseCase,
    getUser *appUser.GetUserUseCase,
    listUsers *appUser.ListUsersUseCase,
    cfg *config.Config,
    log logger.Logger,
) *gin.Engine {
    if cfg.Server.Env == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()
    r.Use(gin.Recovery())

    // Health
    r.GET("/health", HealthHandler)

    // Index 
    r.GET("/", IndexHandler)

  
    attachSwagger(r)

    // Handlers
    authHandler := handlers.NewAuthHandler(requestOTP, verifyOTP, log)
    userHandler := handlers.NewUserHandler(getUser, listUsers, log)

    api := r.Group("/api/v1")
    {
        auth := api.Group("/auth")
        auth.POST("/otp/request", authHandler.RequestOTP)
        auth.POST("/otp/verify", authHandler.VerifyOTP)

        jwtCfg := middleware.JWTConfig{
            Secret:   cfg.JWT.Secret,
            Issuer:   cfg.JWT.Issuer,
            Audience: cfg.JWT.Audience,

            AcceptWithoutBearer: false,
            AcceptFromCookie:    true,
            CookieName:          "access_token",
            AcceptFromQuery:     cfg.Server.Env != "production",
            QueryParam:          "access_token",
        }
        users := api.Group("/users", middleware.AuthMiddleware(jwtCfg))
        users.GET("", userHandler.List)
        users.GET("/:id", userHandler.Get)
    }

    // 404 
    r.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "NOT_FOUND",
            "message": "route not found",
        })
    })

    return r
}

// HealthHandler godoc
// @Summary      Health check
// @Description  Returns service health status
// @Tags         Meta
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func HealthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// IndexHandler godoc
// @Summary      Company landing page
// @Description  Simple HTML landing page for the company
// @Tags         Meta
// @Produce      html
// @Success      200 {string} string "HTML content"
// @Router       / [get]
func IndexHandler(c *gin.Context) {
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
}

// indexHTML contains the HTML for the index page of the OTP Auth Service.
const indexHTML = `<!doctype html>
<html lang="fa" dir="rtl">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>شرکت شما | OTP Auth Service</title>
  <style>
    :root{--bg:#0f172a;--card:#111827;--text:#e5e7eb;--muted:#9ca3af;--primary:#10b981;--accent:#22d3ee}
    *{box-sizing:border-box} body{margin:0;background:radial-gradient(1200px 600px at 80% -10%,#1f2937 0%,#0b1220 55%,#0a0f1c 100%);font-family:IRANSans,Inter,system-ui,Segoe UI,Roboto,Arial,sans-serif;color:var(--text)}
    .wrap{max-width:980px;margin:0 auto;padding:40px 20px}
    header{display:flex;align-items:center;justify-content:space-between;gap:16px}
    .brand{display:flex;align-items:center;gap:12px}
    .logo{width:38px;height:38px;border-radius:8px;background:linear-gradient(135deg,var(--accent),var(--primary));box-shadow:0 10px 30px rgba(16,185,129,.25)}
    .title{font-weight:700;font-size:18px}
    .badge{font-size:12px;color:#0f172a;background:linear-gradient(135deg,#a7f3d0,#67e8f9);padding:4px 10px;border-radius:999px}
    .hero{margin:48px 0 28px}
    .h1{font-size:36px;line-height:1.2;margin:0 0 12px}
    .muted{color:var(--muted)}
    .cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:16px;margin:26px 0 38px}
    .card{background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);backdrop-filter: blur(8px);border-radius:14px;padding:16px 16px 12px}
    .card h3{margin:0 0 8px;font-size:16px}
    .card p{margin:0 0 12px;color:var(--muted);font-size:14px}
    .link{display:inline-block;text-decoration:none;background:linear-gradient(135deg,var(--primary),#34d399);color:#062e2a;padding:8px 12px;border-radius:10px;font-weight:600}
    footer{margin-top:42px;padding-top:16px;border-top:1px dashed rgba(255,255,255,.12);display:flex;justify-content:space-between;font-size:13px;color:var(--muted)}
    @media (max-width:640px){.h1{font-size:28px}}
  </style>
</head>
<body>
  <div class="wrap">
    <header>
      <div class="brand"><div class="logo"></div><div class="title">شرکت شما</div></div>
      <div class="badge">OTP Auth Service</div>
    </header>

    <section class="hero">
      <h1 class="h1">ورود امن با OTP، آماده برای Production</h1>
      <p class="muted">سرویس احراز هویت مبتنی بر OTP با معماری Clean + DDD، مقیاس‌پذیر، امن و مستندسازی‌شده.</p>
    </section>

    <section class="cards">
      <div class="card">
        <h3>سلامت سرویس</h3>
        <p>وضعیت آنی سرویس را بررسی کنید.</p>
        <a class="link" href="/health">نمایش /health</a>
      </div>
      <div class="card">
        <h3>شروع سریع API</h3>
        <p>Endpoints احراز هویت و کاربران.</p>
        <a class="link" href="javascript:void(0)" onclick="alert('POST /api/v1/auth/otp/request\\nPOST /api/v1/auth/otp/verify\\nGET /api/v1/users')">مشاهده</a>
      </div>
      <div class="card">
        <h3>Swagger</h3>
        <p>در صورت فعال‌سازی با build tag.</p>
        <a class="link" href="/swagger/index.html">ورود به Swagger</a>
      </div>
    </section>

    <footer>
      <div>© تمامی حقوق برای شرکت شما محفوظ است.</div>
      <div>نسخه 1.0.0</div>
    </footer>
  </div>
</body>
</html>`