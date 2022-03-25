package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/naixyeur/garagesale/internal/platform/auth"
	"github.com/naixyeur/garagesale/internal/platform/conf"
	"github.com/naixyeur/garagesale/internal/platform/database"
	"github.com/naixyeur/garagesale/internal/schema"
	"github.com/naixyeur/garagesale/internal/user"
	"github.com/pkg/errors"
)

// ================================================================================

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}

}

// ================================================================================
// run
func run() error {

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:true"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	dbConfig := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	var err error
	switch cfg.Args.Num(0) {
	case "migrate":
		err = migrate(dbConfig)
	case "seed":
		err = seed(dbConfig)
	case "useradd":
		err = useradd(dbConfig, cfg.Args.Num(1), cfg.Args.Num(2))
	default:
		err = errors.New("must specify a command")
	}
	if err != nil {
		return err
	}

	return nil
}

// ================================================================================
// migrate
func migrate(dbConfig database.Config) error {
	db, err := database.Open(dbConfig)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := schema.Migrate(db); err != nil {
		return err
	}
	fmt.Println("migration completed")
	return nil
}

// ================================================================================
// seed
func seed(dbConfig database.Config) error {
	db, err := database.Open(dbConfig)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := schema.Seed(db); err != nil {
		return err
	}
	fmt.Println("seeding completed")
	return nil
}

// ================================================================================
// useradd
func useradd(dbConfig database.Config, email, password string) error {
	db, err := database.Open(dbConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	// ==============================

	if email == "" || password == "" {
		return errors.New("useradd command must be called with two additional arguments for email and password")
	}

	fmt.Printf("Admin user will be created with email %q and password %q\n", email, password)
	fmt.Print("Continue? (1/0) ")

	var confirm bool
	if _, err := fmt.Scanf("%t\n", &confirm); err != nil {
		return errors.Wrap(err, "processing response")
	}

	if !confirm {
		fmt.Println("Canceling")
		return nil
	}

	// ==============================

	ctx := context.Background()

	nu := user.NewUser{
		Email:           email,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
		Password:        password,
		PasswordConfirm: password,
	}

	u, err := user.Create(ctx, db, nu, time.Now())
	if err != nil {
		return err
	}

	fmt.Println("User created with id:", u.ID)

	return nil
}
