package rx_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/foreveraloneT/rx"
)

var _ = Describe("Creation", func() {
	Describe("From", func() {
		It("should emit the values from the provided slice", func() {
			slice := []int{1, 5, 3, 4, 3}
			out := From(slice)

			var result []int
			for v := range out {
				result = append(result, v)
			}

			Expect(result).To(Equal(slice))
		})

		Describe("WithBufferSize", func() {
			It("should apply the buffer size option", func() {
				out := From([]string{"hello", "world"}, WithBufferSize(2))

				Expect(cap(out)).To(Equal(2))
			})
		})
	})

	Describe("Interval", func() {
		It("should emit an integer at a fixed interval", func() {
			out := Interval(200 * time.Millisecond)
			result := make([]int, 0)

			ticker := time.NewTicker(700 * time.Millisecond)
			defer ticker.Stop()

		LOOP:
			for {
				select {
				case v := <-out:
					result = append(result, v)
				case <-ticker.C:
					break LOOP
				}
			}

			Expect(result).To(Equal([]int{0, 1, 2}))
		})

		Describe("WithBufferSize", func() {
			It("should apply the buffer size option", func() {
				out := Interval(200*time.Millisecond, WithBufferSize(2))

				Expect(cap(out)).To(Equal(2))
			})
		})
	})
})
