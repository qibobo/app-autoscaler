package db_test

import (
	. "autoscaler/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {
	Describe("DetectDBDriver", func() {
		var (
			dbUrl      string
			driverName string
			err        error
		)
		JustBeforeEach(func() {
			driverName, err = DetectDBDriver(dbUrl)
		})
		Context("when postgres dbUrl is provided", func() {
			BeforeEach(func() {
				dbUrl = "postgres://username:password@localhost/autoscaler?sslmode=disable"
			})
			It("returns PostgresDriverName", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(driverName).To(Equal(PostgresDriverName))
			})
		})
		Context("when mysql dbUrl is provided", func() {
			BeforeEach(func() {
				dbUrl = "mysql://username:password@localhost/autoscaler"
			})
			It("returns MysqlDriverName", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(driverName).To(Equal(MysqlDriverName))
			})
		})
		Context("when other dbUrl is provided", func() {
			BeforeEach(func() {
				dbUrl = "sqlserver://user:pass@remote-host.com/dbname"
			})
			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to detect database type from url"))
			})
		})
	})

})
