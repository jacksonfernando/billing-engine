package main

import (
	"fmt"
	"net/http"
	"os"
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

	// Load configuration from .env file or environment variables
	var configuration global.Configuration

	// Try to read from .env file first
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		// If .env file doesn't exist, use environment variables
		fmt.Println("No .env file found, using environment variables")
	}

	// Set defaults and read from environment variables
	viper.AutomaticEnv()
	viper.SetDefault("host_port", getEnv("APP_PORT", "9006"))
	viper.SetDefault("db_host", getEnv("DB_HOST", "localhost"))
	viper.SetDefault("db_port", getEnv("DB_PORT", "3306"))
	viper.SetDefault("db_name", getEnv("DB_NAME", "billing_engine"))
	viper.SetDefault("db_user", getEnv("DB_USER", "billing_admin"))
	viper.SetDefault("db_pass", getEnv("DB_PASSWORD", "billing_password"))
	viper.SetDefault("private_jwt_access_token_secret", getEnv("JWT_SECRET", "default-secret-key"))
	viper.SetDefault("private_jwt_refresh_token_secret", getEnv("JWT_REFRESH_SECRET", "default-refresh-secret-key"))

	if err := viper.Unmarshal(&configuration); err != nil {
		panic("Unable to decode configuration into struct")
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

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
