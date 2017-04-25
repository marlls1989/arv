package model

import (
	"log"
	"reflect"
)

const Version string = "0.1"

type Model struct {
	start, quit chan struct{}
	memory      *memory
	startPC     uint32
}

/* Since this is a very common element,
 * it is being implemented as a "generic"
 */
func (s *Model) pipeElement(in interface{}, out ...interface{}) {
	vout := reflect.ValueOf(out)
	vin := reflect.ValueOf(in)

	outLen := vout.Len()

	if outLen < 1 {
		log.Panicln("`pipeElement` should contain at least one output argument")
	}

	if vin.Type().Kind() != reflect.Chan {
		log.Panic("`in` must be a channel")
	}

	chType := vin.Type().Elem()

	for i := 0; i < outLen; i++ {
		to := vout.Index(i).Type()
		if (to.Kind() != reflect.Chan) || (to.Elem() == chType) {
			log.Panicf("Argument %d is not a %s channel, but %s",
				i+1, chType.Kind(), to)
		}
	}

	go func() {
		for i := 0; i < outLen; i++ {
			ov := vout.Index(i)
			defer ov.Close()
		}

		for {
			a, ok := vin.Recv()
			if ok {
				for i := 0; i < outLen; i++ {
					ov := vout.Index(i)
					ov.Send(a)
				}
			} else {
				return
			}
		}
	}()
}

func (s *Model) constElement(out chan<- interface{}, val interface{}) {
	outType := reflect.TypeOf(out)
	valType := reflect.TypeOf(val)

	if outType.Elem() != valType {
		log.Panicf("`out` channel type must match `val` type %s, but got %s",
			valType.Kind(), outType.Elem().Kind())
	}

	go func() {
		defer close(out)
		for {
			select {
			case <-s.quit:
				return
			case out <- val:
			}
		}
	}()
}
