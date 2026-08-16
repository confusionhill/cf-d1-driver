module github.com/confusionhill/cf-d1-driver/examples/sqlx_app

go 1.22

require (
	github.com/confusionhill/cf-d1-driver v0.0.1
	github.com/jmoiron/sqlx v1.4.0
)

replace github.com/confusionhill/cf-d1-driver => ../..
