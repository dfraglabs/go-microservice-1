package errtype

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SampleErr struct {
	e error
}
func (t SampleErr) Error() string { return t.e.Error() }

var _ = Describe("errtype", func() {
	var (
		t string
		e error
	)

	JustBeforeEach(func() {
		switch e.(type) {
		case TypedErr:
			t = "Typed"
		case SampleErr:
			t = "Sample"
		case APINotFoundErr:
			t = "NotFound"
		default:
			t = "other"
		}
	})

	Context("real typed error", func() {
		BeforeEach(func() {
			e = func() error { return APINotFoundErr{errors.New("user not found")} }()
		})

		It("type can be detected", func() {
			Expect(t).To(Equal("NotFound"))
		})
	})

	Context("base typed error", func() {
		BeforeEach(func() {
			e = func() error { return TypedErr{errors.New("user not found")} }()
		})

		It("type can be detected", func() {
			Expect(t).To(Equal("Typed"))
		})
	})

	Context("independently declared", func() {
		BeforeEach(func() {
			e = func() error { return SampleErr{errors.New("user not found")} }()
		})

		It("type can be detected", func() {
			Expect(t).To(Equal("Sample"))
		})
	})

	Context("other error", func() {
		BeforeEach(func() {
			e = func() error { return errors.New("user not found") }()
		})

		It("type can be detected", func() {
			Expect(t).To(Equal("other"))
		})
	})
})
