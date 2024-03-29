package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/0xThomas3000/food_delivery/components/appctx"
	"github.com/0xThomas3000/food_delivery/components/uploadprovider"
	"github.com/0xThomas3000/food_delivery/middleware"
	"github.com/0xThomas3000/food_delivery/pubsub/localpb"
	"github.com/0xThomas3000/food_delivery/routes"
	"github.com/0xThomas3000/food_delivery/subscriber"
	"github.com/0xThomas3000/food_delivery/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	s3Provider := uploadprovider.NewS3Provider(
		config.S3BucketName,
		config.S3Region,
		config.S3APIKey,
		config.S3SecretKey,
		config.S3Domain,
	)

	secretKey := config.SecretKey
	ps := localpb.NewPubSub()
	appContext := appctx.NewAppContext(db, s3Provider, secretKey, ps)

	// Setup subscribers
	// subscriber.Setup(appContext, context.Background())
	_ = subscriber.NewEngine(appContext).Start()

	r := gin.Default()
	r.Use(middleware.Recover(appContext))

	// Đăng ký link cho cái static để hiển thị hình
	r.Static("/static", "./static") // Đi search mục "static" => gin sẽ kiếm thư mục "static" để đọc

	v1 := r.Group("/v1")

	routes.SetupRoute(appContext, v1)
	routes.SetupAdminRoute(appContext, v1)

	r.Run()
}
