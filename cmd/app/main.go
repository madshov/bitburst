package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	repo "github.com/madshov/bitburst/app/repository/postgres"
	service "github.com/madshov/bitburst/app/service"
	httptr "github.com/madshov/bitburst/app/transport"
)

func main() {
	logger := log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime)

	logger.Println("starting test service")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"bitburst_postgresql_1", 5432, "bitburst", "bitburst", "bitburst")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatalln("msg", "db connection failed", "reason", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatalln("msg", "db connection failed", "reason", err)
	}

	repo := repo.NewObjectRepo(db)
	svc := service.NewService(logger, repo)
	hand := httptr.NewObjectHandler(logger, svc)

	http.Handle("/callback", hand)

	errs := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":9090")
		if err != nil {
			errs <- err
			return
		}

		errs <- http.Serve(listener, nil)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Println("msg", "app terminated", "reason", <-errs)
}
