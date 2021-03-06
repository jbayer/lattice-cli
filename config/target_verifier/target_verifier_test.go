package target_verifier_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/cloudfoundry-incubator/receptor/fake_receptor"

	"github.com/pivotal-cf-experimental/lattice-cli/config/target_verifier"
)

type fakereceptorClientBuilder struct {
	receptorClient receptor.Client
}

func (f *fakereceptorClientBuilder) Build(target string) receptor.Client {
	return f.receptorClient
}

var _ = Describe("targetVerifier", func() {
	Describe("ValidateAuthorization", func() {
		var fakeReceptorClient *fake_receptor.FakeClient
		var targets []string

		var fakeReceptorClientFactory = func(target string) receptor.Client {
			targets = append(targets, target)
			return fakeReceptorClient
		}

		BeforeEach(func() {
			fakeReceptorClient = &fake_receptor.FakeClient{}
			targets = []string{}
		})

		It("returns true if the receptor does not return an error", func() {
			fakeReceptorClient.DesiredLRPsReturns([]receptor.DesiredLRPResponse{}, nil)
			targetVerifier := target_verifier.New(fakeReceptorClientFactory)

			ok, err := targetVerifier.ValidateAuthorization("http://receptor.mylattice.com")
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
			Expect(targets).To(Equal([]string{"http://receptor.mylattice.com"}))
		})

		It("returns false if the receptor returns an authorization error", func() {
			fakeReceptorClient.DesiredLRPsReturns([]receptor.DesiredLRPResponse{}, errors.New(receptor.Unauthorized))
			targetVerifier := target_verifier.New(fakeReceptorClientFactory)

			ok, err := targetVerifier.ValidateAuthorization("http://receptor.mylattice.com")
			Expect(ok).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if the receptor returns an error other than unauthorized", func() {
			fakeReceptorClient.DesiredLRPsReturns([]receptor.DesiredLRPResponse{}, errors.New(receptor.UnknownError))
			targetVerifier := target_verifier.New(fakeReceptorClientFactory)

			_, err := targetVerifier.ValidateAuthorization("http://receptor.mylattice.com")
			Expect(err).To(HaveOccurred())
		})
	})
})
