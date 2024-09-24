package internal

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestRecover(t *testing.T) {
	var output bytes.Buffer

	var logger *log.Logger = log.New(&output, "", 0)

	var yolk = func(recover *Recover) {
		logger.Print("[yolk] begin")

		recover.
			Defer(func(err interface{}) {
				logger.Print("[yolk] defer")
			}).
			Do(func(f Finalizer) {
				logger.Print("[yolk] do")

				f.Add(func(err interface{}) {
					logger.Print("[yolk] Finalizer")
				})
			})

		logger.Print("[yolk] end")
	}

	var albumen = func(recover *Recover) {
		logger.Print("[albumen] begin")

		recover.
			Defer(func(err interface{}) {
				logger.Print("[albumen] defer")
			}).
			Do(func(f Finalizer) {
				logger.Print("[albumen] do")

				f.Add(func(err interface{}) {
					logger.Print("[albumen] Finalizer")
				})
				yolk(recover)
			})

		logger.Print("[albumen] end")
	}

	var shell = func(recover *Recover) {
		logger.Print("[shell] begin")

		recover.
			Defer(func(err interface{}) {
				logger.Print("[shell] defer")
			}).
			Do(func(f Finalizer) {
				logger.Print("[shell] do")

				f.Add(func(err interface{}) {
					logger.Print("[shell] Finalizer")
				})
				albumen(recover)
			})

		logger.Print("[shell] end")
	}

	func() {
		logger.Print("begin")

		defer func() {
			logger.Print("end")
		}()

		shell(new(Recover))
	}()

	expectedOutput := strings.Join([]string{
		"begin\n",
		"[shell] begin\n",
		"[shell] do\n",
		"[albumen] begin\n",
		"[albumen] do\n",
		"[yolk] begin\n",
		"[yolk] do\n",
		"[yolk] defer\n",
		"[yolk] Finalizer\n",
		"[yolk] end\n",
		"[albumen] defer\n",
		"[albumen] Finalizer\n",
		"[albumen] end\n",
		"[shell] defer\n",
		"[shell] Finalizer\n",
		"[shell] end\n",
		"end\n",
	}, "")
	if expectedOutput != output.String() {
		t.Errorf("assert output:: expected '%v', got '%v'", expectedOutput, output.String())
	}
}
