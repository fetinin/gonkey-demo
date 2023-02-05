package gonkey__test

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lamoda/gonkey/fixtures"
	"github.com/lamoda/gonkey/mocks"
	"github.com/lamoda/gonkey/runner"
	"github.com/stretchr/testify/require"

	app "gonkey-example/case-app/internal"
)

const dbDSN = "postgres://service:service@localhost:6543/service?sslmode=disable"

func TestFuncCases(t *testing.T) {
	db, err := app.NewDB(context.Background(), dbDSN)
	require.NoError(t, err)

	m := mocks.NewNop("nameApi")
	err = m.Start()
	require.NoError(t, err)

	nicksGenAddr := "http://" + m.Service("nameApi").ServerAddr()
	api := app.NewAPI(db, nicksGenAddr)

	srv := httptest.NewServer(api)

	// run test cases from your dir with Allure report generation
	runner.RunWithTesting(t, &runner.RunWithTestingParams{
		Server:      srv,
		TestsDir:    "cases",
		FixturesDir: "fixtures",
		DbType:      fixtures.Postgres,
		DB:          newTestDBConn(t, dbDSN),
		Mocks:       m,
	})
}

func newTestDBConn(t *testing.T, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)

	require.NoError(t, db.Ping())
	return db
}
