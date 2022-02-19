package repository_test

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/sidelight-labs/libdb/repository"
	"os"
	"testing"
)

func TestUnitRepository(t *testing.T) {
	spec.Run(t, "Repository", testRepository, spec.Report(report.Terminal{}))
}

func testRepository(t *testing.T, when spec.G, it spec.S) {
	var (
		subject repository.Store
	)

	it.Before(func() {
		RegisterTestingT(t)
	})

	it.After(func() {
		subject.Close()
	})

	when("NewMySQL()", func() {
		var databaseName = "database_name"

		it("returns an error when no hostname is specified", func() {
			Expect(os.Unsetenv(repository.DBHostEnv)).To(Succeed())
			var err error
			subject, err = repository.NewMySQL(databaseName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(repository.DatabaseErrorPrefix))
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(repository.MissingEnvError, repository.DBHostEnv)))
		})
		it("returns an error when the port value is not an integer", func() {
			Expect(os.Setenv(repository.DBHostEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBUserEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPasswordEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPortEnv, "test")).To(Succeed())
			var err error
			subject, err = repository.NewMySQL(databaseName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(repository.PortError))
		})
		it("returns NO error if all environment variables are valid values", func() {
			Expect(os.Setenv(repository.DBHostEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBUserEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPasswordEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPortEnv, "123")).To(Succeed())
			var err error
			subject, err = repository.NewMySQL(databaseName)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	when("NewPostgres()", func() {
		var databaseName = "database_name"

		it("returns an error when no hostname is specified", func() {
			Expect(os.Unsetenv(repository.DBHostEnv)).To(Succeed())
			var err error
			subject, err = repository.NewPostgres(databaseName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(repository.DatabaseErrorPrefix))
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(repository.MissingEnvError, repository.DBHostEnv)))
		})
		it("returns an error when the port value is not an integer", func() {
			Expect(os.Setenv(repository.DBHostEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBUserEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPasswordEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPortEnv, "test")).To(Succeed())
			var err error
			subject, err = repository.NewPostgres(databaseName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(repository.PortError))
		})
		it("returns NO error if all environment variables are valid values", func() {
			Expect(os.Setenv(repository.DBHostEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBUserEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPasswordEnv, "test")).To(Succeed())
			Expect(os.Setenv(repository.DBPortEnv, "123")).To(Succeed())
			var err error
			subject, err = repository.NewPostgres(databaseName)
			Expect(err).NotTo(HaveOccurred())
		})
	})
}
