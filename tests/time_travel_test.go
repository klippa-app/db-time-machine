package tests_test

import (
	"context"

	"github.com/klippa-app/db-time-machine/db"
	"github.com/klippa-app/db-time-machine/db/mocks"
	"github.com/klippa-app/db-time-machine/internal"
	"github.com/klippa-app/db-time-machine/internal/config"
	"github.com/klippa-app/db-time-machine/internal/hashes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("TimeTravel", func() {
	var ctx context.Context
	var cfg config.Config

	BeforeEach(func() {
		ctx = context.Background()
		cfg = config.Config{
			Prefix: "test",
			Connection: config.ConnectionConfig{
				URI:      "someTestUri",
				Database: "testDatabase",
			},
			Migration: config.MigrationConfig{
				Directory: ".",
				Format:    ".",
				Command:   ".",
			},
		}
		ctx = config.Attach(ctx, cfg)
	})

	Context("database has current migration", func() {
		var tempCtx context.Context
		var mockDriver *mocks.Driver

		BeforeEach(func() {
			hashList := []string{
				"098f6bcd4621d373cade4e832627b4f6",
				"925e267e725330ded9fe700cdf8669e2",
				"4984f75137ddea54b48bd574c5fc1000",
				"6c2f4616f4b3daf1fe5edbe53effbc63",
			}

			tempCtx = hashes.Attach(ctx, hashList)

			mockDriver = new(mocks.Driver)
			mockDriver.On("List", mock.Anything).Return([]string{"test_925e267e", "test_4984f751", "test_6c2f4616", "test_098f6bcd"}, nil)
			tempCtx = db.AttachContext(tempCtx, mockDriver)
		})

		It("will return the existing database name", func() {
			result, err := internal.TimeTravel(tempCtx,
				func(ctx context.Context, target string) error {
					return nil
				})
			Expect(result).To(BeEquivalentTo("test_098f6bcd"))
			Expect(err).ToNot(HaveOccurred())

			Expect(mockDriver.AssertCalled((mock.TestingT)(GinkgoT()), "List", mock.Anything)).To(BeTrue())
			Expect(mockDriver.AssertNotCalled((mock.TestingT)(GinkgoT()), "Clone", mock.Anything, mock.Anything, mock.Anything)).To(BeTrue())
			Expect(mockDriver.AssertNotCalled((mock.TestingT)(GinkgoT()), "Create", mock.Anything, mock.Anything)).To(BeTrue())
		})
	})

	Context("database is one migration behind", func() {
		var tempCtx context.Context
		var mockDriver *mocks.Driver

		BeforeEach(func() {
			hashList := []string{
				"098f6bcd4621d373cade4e832627b4f6",
				"925e267e725330ded9fe700cdf8669e2",
				"4984f75137ddea54b48bd574c5fc1000",
				"6c2f4616f4b3daf1fe5edbe53effbc63",
			}

			tempCtx = hashes.Attach(ctx, hashList)

			mockDriver = new(mocks.Driver)
			mockDriver.On("List", mock.Anything).Return([]string{"test_925e267e", "test_4984f751", "test_6c2f4616"}, nil)
			mockDriver.On("Clone", mock.Anything, "test_925e267e", "test_098f6bcd").Return(nil)
			tempCtx = db.AttachContext(tempCtx, mockDriver)
		})

		It("will return a new cloned database", func() {
			result, err := internal.TimeTravel(tempCtx, func(ctx context.Context, target string) error {
				return nil
			})

			Expect(result).To(BeEquivalentTo("test_098f6bcd"))
			Expect(err).ToNot(HaveOccurred())

			Expect(mockDriver.AssertCalled((mock.TestingT)(GinkgoT()), "List", mock.Anything)).To(BeTrue())
			Expect(mockDriver.AssertCalled((mock.TestingT)(GinkgoT()), "Clone", mock.Anything, mock.Anything, mock.Anything)).To(BeTrue())
			Expect(mockDriver.AssertNotCalled((mock.TestingT)(GinkgoT()), "Create", mock.Anything, mock.Anything)).To(BeTrue())
			Expect(mockDriver.AssertNotCalled((mock.TestingT)(GinkgoT()), "Remove", mock.Anything, mock.Anything)).To(BeTrue())
		})
	})

	Context("database has no common migrations", func() {
		var tempCtx context.Context
		var mockDriver *mocks.Driver

		BeforeEach(func() {
			hashList := []string{
				"098f6bcd4621d373cade4e832627b4f6",
				"925e267e725330ded9fe700cdf8669e2",
				"4984f75137ddea54b48bd574c5fc1000",
				"6c2f4616f4b3daf1fe5edbe53effbc63",
			}

			tempCtx = hashes.Attach(ctx, hashList)

			mockDriver = new(mocks.Driver)
			mockDriver.On("List", mock.Anything).Return([]string{"test_8725e267e", "test_3984f751", "test_5c2f4616"}, nil)
			mockDriver.On("Create", mock.Anything, "test_098f6bcd").Return(nil)
			tempCtx = db.AttachContext(tempCtx, mockDriver)
		})

		It("will return a new database that is freshly created not cloned", func() {
			result, err := internal.TimeTravel(tempCtx, func(ctx context.Context, target string) error {
				return nil
			})

			Expect(result).To(BeEquivalentTo("test_098f6bcd"))
			Expect(err).ToNot(HaveOccurred())

			Expect(mockDriver.AssertCalled((mock.TestingT)(GinkgoT()), "List", mock.Anything)).To(BeTrue())
			Expect(mockDriver.AssertCalled((mock.TestingT)(GinkgoT()), "Create", mock.Anything, "test_098f6bcd")).To(BeTrue())
			Expect(mockDriver.AssertNotCalled((mock.TestingT)(GinkgoT()), "Clone", mock.Anything, mock.Anything, mock.Anything)).To(BeTrue())
			Expect(mockDriver.AssertNotCalled((mock.TestingT)(GinkgoT()), "Remove", mock.Anything, mock.Anything)).To(BeTrue())
		})
	})
})
