package main

import (
	"context"
	"log"
	"os"

	"ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("migration name is required")
	}

	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		configPath = "config.yaml"
	}

	dir, err := migrate.NewLocalDir("./ent/migrations")
	if err != nil {
		log.Fatalln(err)
	}
	graph, err := entc.LoadGraph("./ent/schema", &gen.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	tbls, err := graph.Tables()
	if err != nil {
		log.Fatalln(err)
	}
	cfg := harvester.GetConfig(configPath, true)
	mysqlConfig := mysql.GetMYSQLConfig(&cfg.MySQL)
	drv, err := sql.Open("mysql", mysqlConfig.FormatDSN())

	if err != nil {
		log.Fatalln(err)
	}

	m, err := schema.NewMigrate(drv, schema.WithDir(dir))
	if err != nil {
		log.Fatalln(err)
	}
	if err := m.NamedDiff(context.Background(), os.Args[1], tbls...); err != nil {
		log.Fatalln(err)
	}
}
