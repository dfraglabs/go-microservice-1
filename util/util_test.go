package util

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util", func() {

	Describe("RandString", func() {
		const (
			requestLen = 32
		)

		Context("called once", func() {
			var (
				returnedString string
			)

			JustBeforeEach(func() {
				returnedString = RandString(requestLen, time.Now())
			})
			It("should have correct len", func() {
				Expect(len(returnedString)).To(Equal(requestLen))
			})
		})

		Context("different times passed", func() {
			var (
				string1 string
				string2 string
			)

			JustBeforeEach(func() {
				string1 = RandString(requestLen, time.Unix(1257894000, 0))
				string2 = RandString(requestLen, time.Unix(1257895000, 0))
			})

			It("should return different strings", func() {
				Expect(string1).ToNot(Equal(string2))
			})
		})

		Context("same times passed", func() {
			var (
				string1 string
				string2 string
			)

			JustBeforeEach(func() {
				string1 = RandString(requestLen, time.Unix(1257894000, 0))
				string2 = RandString(requestLen, time.Unix(1257894000, 0))
			})

			It("should return different strings", func() {
				Expect(string1).To(Equal(string2))
			})
		})
	})

})
