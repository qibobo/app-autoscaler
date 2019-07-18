package custom_metrics_cred_helper_test

import (
	"errors"

	. "autoscaler/api/custom_metrics_cred_helper"
	"autoscaler/fakes"
	"autoscaler/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CustomMetricCredHelper", func() {
	var (
		policyDB   *fakes.FakePolicyDB
		appId      = "testAppId"
		credential = &models.CustomMetricCredentials{
			Username: "username",
			Password: "password",
		}
	)
	BeforeEach(func() {
		policyDB = &fakes.FakePolicyDB{}
	})
	Context("CreateCustomMetricsCredential", func() {
		var err error
		JustBeforeEach(func() {
			err = CreateCustomMetricsCredential(appId, policyDB, MaxRetry)
		})
		Context("when there is no error when calling policydb", func() {
			BeforeEach(func() {
				policyDB.SaveCustomMetricsCredReturns(nil)
			})
			It("should try only once and succeed", func() {
				Expect(policyDB.SaveCustomMetricsCredCallCount()).To(Equal(1))
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when there is continous error when calling policydb", func() {
			BeforeEach(func() {
				policyDB.SaveCustomMetricsCredReturns(errors.New("dberror"))
			})
			It("should try MaxRetry times and return error", func() {
				Expect(policyDB.SaveCustomMetricsCredCallCount()).To(Equal(MaxRetry))
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Context("DeleteCustomMetricsCredential", func() {
		var err error
		JustBeforeEach(func() {
			err = DeleteCustomMetricsCredential(appId, policyDB, MaxRetry)
		})
		Context("when there is no error when calling policydb", func() {
			BeforeEach(func() {
				policyDB.DeleteCustomMetricsCredReturns(nil)
			})
			It("should try only once and succeed", func() {
				Expect(policyDB.DeleteCustomMetricsCredCallCount()).To(Equal(1))
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when there is continous error when calling policydb", func() {
			BeforeEach(func() {
				policyDB.DeleteCustomMetricsCredReturns(errors.New("dberror"))
			})
			It("should try MaxRetry times and return error", func() {
				Expect(policyDB.DeleteCustomMetricsCredCallCount()).To(Equal(MaxRetry))
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Context("GetCustomMetricsCredential", func() {
		var err error
		var cred *models.CustomMetricCredentials
		JustBeforeEach(func() {
			cred, err = GetCustomMetricsCredential(appId, policyDB, MaxRetry)
		})
		Context("when there is no error when calling policydb", func() {
			BeforeEach(func() {
				policyDB.GetCustomMetricsCredsReturns(credential, nil)
			})
			It("should try only once and succeed", func() {
				Expect(policyDB.GetCustomMetricsCredsCallCount()).To(Equal(1))
				Expect(err).NotTo(HaveOccurred())
				Expect(cred).To(Equal(credential))
			})
		})
		Context("when there is continous error when calling policydb", func() {
			BeforeEach(func() {
				policyDB.GetCustomMetricsCredsReturns(nil, errors.New("dberror"))
			})
			It("should try MaxRetry times and return error", func() {
				Expect(policyDB.GetCustomMetricsCredsCallCount()).To(Equal(MaxRetry))
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
