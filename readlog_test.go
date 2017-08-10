package ccruncher_test

import (
	. "ccruncher"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Readlog", func() {
	var (
		ccLog *CCLog
		log   *os.File
		err   error
	)
	// var logFile *bytes.Reader

	const log1 = "fixtures/cclog.1"
	const log2 = "fixtures/cclog.2"
	const log3 = "fixtures/cclog.3"
	const log4 = "fixtures/cclog.4"

	BeforeEach(func() {
		log, err = os.Open(log2)
		Expect(err).NotTo(HaveOccurred())

		defer log.Close()
		ccLog, err = ParseLog(log)

		Expect(err).NotTo(HaveOccurred())
	})

	Describe("EntriesForRequest", func() {
		var requestID string
		var entries []LogEntry

		BeforeEach(func() {
			requestID = "d462237c-c201-4ab2-57c3-e1e79570b0b0::0bfa387c-38df-4053-9dc9-3205d102ddf2"
			entries = ccLog.EntriesForRequest(requestID)
		})

		It("Returns the entries for the specified request", func() {
			Expect(entries).To(HaveLen(6))
		})

		It("Links each entry to the original request", func() {
			for _, entry := range entries {
				Expect(entry.RequestID()).To(Equal("d462237c-c201-4ab2-57c3-e1e79570b0b0::0bfa387c-38df-4053-9dc9-3205d102ddf2"))
				Expect(entry.AppGUID()).To(Equal("3b664a4a-3aba-4a1b-baa7-8ab67b4ade96"))
				Expect(entry.HttpMethod()).To(Equal("GET"))
				Expect(entry.URIPath()).To(Equal("/v2/apps/3b664a4a-3aba-4a1b-baa7-8ab67b4ade96/stats"))
			}
		})

		It("Does Stuff", func() {
			a, _ := entries[0].Render()
			Expect(string(a)).To(Equal("a"))
		})
	})

	Describe("RequestsForApp", func() {
		var appGUID string

		BeforeEach(func() {
			appGUID = "1b86d700-1c7f-4479-b232-d677e85131b2"
		})

		It("Returns the ids of the requests for the specified app", func() {
			ids := ccLog.RequestsForApp(appGUID)
			Expect(ids).To(HaveLen(2))
		})
	})

	// Describe("LogEntry", func() {
	// 	Describe("String()", func() {
	// 		It("Returns the correct")
	// 	})
	// })

	Describe("Apps", func() {
		BeforeEach(func() {
			log, err = os.Open(log4)
			Expect(err).NotTo(HaveOccurred())

			defer log.Close()
			ccLog, err = ParseLog(log)

			Expect(err).NotTo(HaveOccurred())
		})
		It("Returns a list of the app guids in the log", func() {
			Expect(ccLog.Apps()).To(HaveLen(6))
		})
	})
	Describe("Entries()", func() {
		BeforeEach(func() {

			ccLog = &CCLog{}
		})

		It("Returns all of the log entries", func() {

		})
	})

	Describe("ParseLog", func() {
		var ()

		It("Parses the log", func() {

			log, err = os.Open(log2)
			Expect(err).NotTo(HaveOccurred())

			defer log.Close()
			ccLog, err = ParseLog(log)

			Expect(err).NotTo(HaveOccurred())
			entries := ccLog.Entries()

			Expect(len(entries)).To(Equal(300))
		})

		Context("When the log is not valid", func() {
			It("Returns an error", func() {
				log, err = os.Open("fixtures/badlog")
				Expect(err).NotTo(HaveOccurred())

				defer log.Close()
				_, err := ParseLog(log)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
