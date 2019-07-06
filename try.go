package ex

func Try(exec func()) (context *tryContext) {
	context = &tryContext{}
	defer func() {
		if e := recover(); e != nil {
			context.Error = e
		}
	}()
	exec()
	return context
}

type tryContext struct {
	Error interface{}
}

func (t *tryContext) Catch(exec func(err interface{})) *tryContext {
	if t.Error != nil {
		exec(t.Error)
	}
	return t
}

func (t *tryContext) SafeCatch(exec func(err interface{})) *tryContext {
	if t.Error != nil {
		Try(func() {
			exec(t.Error)
		}).Catch(func(err interface{}) {
			Wrap(err).PrintErrorStack()
		})
	}
	return t
}

func (t *tryContext) OK() bool {
	return t.Error == nil
}
