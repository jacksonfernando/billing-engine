package main

import (
	"fmt"
	"net/http"
	"time"

	disbursementHTTPHandler "billing-engine/disbursement/handler/http"
	disbursementRepository "billing-engine/disbursement/repository/mysql"
	disbursementService "billing-engine/disbursement/service"

	repaymentHTTPHandler "billing-engine/repayment/handler/http"
	repaymentRepository "billing-engine/repayment/repository/mysql"
	repaymentService "billing-engine/repayment/service"

	loanQueryHTTPHandler "billing-engine/loan_query/handler/http"
	loanQueryRepository "billing-engine/loan_query/repository/mysql"
	loanQueryService "billing-engine/loan_query/service"

	"billing-engine/global"
	"billing-engine/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic("Unable to load Asia/Jakarta timezone")
	}
	time.Local = loc

	viper.SetConfigFile(".env")
	var configuration global.Configuration
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read .env file")
	}
	if err := viper.Unmarshal(&configuration); err != nil {
		panic("Unable to decode into struct")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configuration.DbUser,
		configuration.DbPass,
		configuration.DbHost,
		configuration.DbPort,
		configuration.DbName,
	)
	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database")
	}

	// initialize echo web apps
	newEcho := echo.New()
	middlewares := middlewares.InitMiddleware([]byte(configuration.PrivateJWTAccessTokenSecret))
	newEcho.Use(middlewares.ValidateCORS)

	newEcho.GET("/health", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, map[string]interface{}{"message": "Billing Engine is live"})
	})

	// Initialize disbursement module
	disbursementRepo := disbursementRepository.NewDisbursementMySQLRepository(mysqlDb)
	disbursementSvc := disbursementService.NewDisbursementService(disbursementRepo)
	disbursementHTTPHandler.NewDisbursementHandler(newEcho, disbursementSvc, middlewares)

	// Initialize repayment module
	repaymentRepo := repaymentRepository.NewRepaymentMySQLRepository(mysqlDb)
	repaymentSvc := repaymentService.NewRepaymentService(repaymentRepo)
	repaymentHTTPHandler.NewRepaymentHandler(newEcho, repaymentSvc, middlewares)

	// Initialize loan query module
	loanQueryRepo := loanQueryRepository.NewLoanQueryMySQLRepository(mysqlDb)
	loanQuerySvc := loanQueryService.NewLoanQueryService(loanQueryRepo)
	loanQueryHTTPHandler.NewLoanQueryHandler(newEcho, loanQuerySvc, middlewares)

	newEcho.Logger.Fatal(newEcho.Start(fmt.Sprintf(":%s", configuration.HostPort)))
}
