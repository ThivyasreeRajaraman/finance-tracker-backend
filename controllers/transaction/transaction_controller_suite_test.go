package transactioncontrollers

import (
	"os"
	"testing"

	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	// initializers.LoadEnvVariables(true)
	initializers.InitDB(os.Getenv("TEST_DATABASE_URL"))
	initializers.SyncDatabase()
}

func TestControllerService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	initializers.DB.Exec("TRUNCATE TABLE users CASCADE")
	initializers.DB.Exec("TRUNCATE TABLE categories CASCADE")
	initializers.DB.Exec("TRUNCATE TABLE transaction_partners CASCADE")
	initializers.DB.Exec("TRUNCATE TABLE transactions CASCADE")
})
