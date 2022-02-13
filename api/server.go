package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/unrolled/secure"
	"goddns/api/controllers"
	_ "goddns/api/docs"
	"goddns/api/middlewares"
	"goddns/bolts/store"
	"goddns/config"
	"goddns/metrics"
)

type Server struct {
	conf      *config.API
	datastore *store.Store
	router    *gin.Engine
}

// Swagger API document annotations
// @title GoDDNS API
// @version 1.0
// @description Use authorized client to config GoDDNS.
// @host localhost:8001
// @BasePath /api/v1
func NewServer(conf *config.API, store *store.Store, metrics metrics.Metrics) (*Server, error) {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	if metrics != nil {
		router.Use(middlewares.Metrics(metrics))
	}
	if !conf.AllowInsecureHTTP {
		router.Use(middlewares.TlsHandler(secure.Options{
			SSLRedirect: true,
		}))
	}
	if conf.AllowCors {
		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			AllowFiles:       true,
			AllowWebSockets:  true,
		}))
	}
	if conf.Swagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Route UI resource if assigned
	// Route to index.html in backend when using Vue history mode router,
	// see: https://github.com/gin-gonic/gin/issues/1048
	if conf.UIResource != "" {
		router.Use(static.Serve(conf.UIPath, static.LocalFile(conf.UIResource, true)))
		router.NoRoute(func(c *gin.Context) {
			c.File(conf.UIResource + "/index.html")
		})
	}

	// login route
	router.POST("/login", controllers.Authenticate(store))

	// authenticated routes which require login
	sv1Admin := router.Group("/api/v1/")
	sv1Admin.Use(middlewares.JWTAuth())
	{
		sv1Admin.POST("/domain", controllers.AddDomain(store))
		sv1Admin.DELETE("/domain", controllers.DeleteDomain(store))
		sv1Admin.GET("/domain/provider/list", controllers.ListDomainProviders(store))
		sv1Admin.GET("/domain/list", controllers.ListDomains(store))

		sv1Admin.POST("/record", controllers.AddRecord(store))
		sv1Admin.DELETE("/record", controllers.DeleteRecord(store))
		sv1Admin.GET("/record/list", controllers.ListRecords(store))

		sv1Admin.GET("/ip/current", controllers.LookupCurrentIP(store))
		sv1Admin.GET("/ip/last", controllers.GetLastIP(store))

		sv1Admin.PUT("/ddns/execute", controllers.DDNSExecute(store))

		sv1Admin.GET("/user/current", controllers.ParseUser(store))
	}

	return &Server{
		conf:      conf,
		datastore: store,
		router:    router,
	}, nil
}

// Serve API, this is a block method
func (s *Server) Serve() error {
	if s.conf.AllowInsecureHTTP {
		return s.router.Run(s.conf.Listen)
	}
	return s.router.RunTLS(s.conf.Listen, s.conf.TLSCert, s.conf.TLSKey)
}
