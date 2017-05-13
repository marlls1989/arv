package processor

import (
	"bitbucket.org/marcos_sartori/qdi-riscv/memory"
	"log"
	"reflect"
)

type Processor struct {
	start, quit chan struct{}
	Memory      memory.Memory
	regFile     regFile
	startPC     uint32
}

func encodeOneHot32(val ...uint) (ret uint32) {
	ret = 0
	for _, v := range val {
		ret |= 1 << (v & 0x1F)
	}

	return
}

func decodeOneHot32(val uint32) (ret []uint) {
	var i uint
	for i = 0; i < 32; i++ {
		if (val & 1) != 0 {
			ret = append(ret, i)
		}

		val >>= 1
	}

	return
}

/* Since this is a very common element,
 * it is being implemented as a "generic"
 */
func (s *Processor) pipeElement(in interface{}, out ...interface{}) {
	vout := reflect.ValueOf(out)
	vin := reflect.ValueOf(in)

	outLen := vout.Len()

	if outLen < 1 {
		log.Panicln("`pipeElement` should contain at least one output argument")
	}

	inputIsChannel := false
	chType := vin.Type()

	if chType.Kind() == reflect.Chan {
		chType = chType.Elem()
		inputIsChannel = true
	}

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

		if inputIsChannel {
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
		} else {
			var cases []reflect.SelectCase

			cases[0] = reflect.SelectCase{
				Chan: reflect.ValueOf(s.quit),
				Dir:  reflect.SelectRecv}

			for i := 0; i < outLen; i++ {
				cases[i+1] = reflect.SelectCase{
					Chan: vout.Index(i),
					Dir:  reflect.SelectSend,
					Send: vin}
			}

			for {
				_, _, ok := reflect.Select(cases)
				if !ok {
					return
				}
			}
		}
	}()
}

func (s *Processor) pipeElementWithInitization(in interface{}, init interface{}, out ...interface{}) {
	vout := reflect.ValueOf(out)
	vin := reflect.ValueOf(in)
	vinit := reflect.ValueOf(init)

	outLen := vout.Len()

	if outLen < 1 {
		log.Panic("`pipeElement` should contain at least one output argument")
	}

	inputIsChannel := false
	chType := vin.Type()

	if chType.Kind() == reflect.Chan {
		chType = chType.Elem()
		inputIsChannel = true
	}

	if vinit.Type() != chType {
		log.Panic("`pipeElement` initialization must be the same type as main input and output")
	}

	for i := 0; i < outLen; i++ {
		to := vout.Index(i).Type()
		if (to.Kind() != reflect.Chan) || (to.Elem() != chType) {
			log.Panicf("Output %d is not a %s channel, but %s",
				i, chType.Kind(), to)
		}
	}

	go func() {
		for i := 0; i < outLen; i++ {
			ov := vout.Index(i)
			defer ov.Close()
			ov.Send(vinit)
		}

		if inputIsChannel {
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
		} else {
			var cases []reflect.SelectCase

			cases[0] = reflect.SelectCase{
				Chan: reflect.ValueOf(s.quit),
				Dir:  reflect.SelectRecv}

			for i := 0; i < outLen; i++ {
				cases[i+1] = reflect.SelectCase{
					Chan: vout.Index(i),
					Dir:  reflect.SelectSend,
					Send: vin}
			}

			for {
				_, _, ok := reflect.Select(cases)
				if !ok {
					return
				}
			}
		}
	}()
}

func (s *Processor) Start() {
	if s.Memory != nil {
		close(s.start)
	} else {
		log.Panic("Processor has no memory attached")
	}
}

func (s *Processor) Stop() {
	close(s.quit)
}
