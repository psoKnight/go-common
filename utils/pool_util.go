package utils

type GoPool struct {
	MaxLimit  int
	tokenChan chan struct{}
}

type GoPoolOption func(*GoPool)

func WithMaxLimit(max int) GoPoolOption {
	return func(gp *GoPool) {
		gp.MaxLimit = max
		gp.tokenChan = make(chan struct{}, gp.MaxLimit)

		for i := 0; i < gp.MaxLimit; i++ {
			gp.tokenChan <- struct{}{}
		}
	}
}

func NewGoPool(options ...GoPoolOption) *GoPool {
	p := &GoPool{}
	for _, o := range options {
		o(p)
	}

	return p
}

// Submit 提交会等待一个token 令牌，然后执行fn
func (gp *GoPool) Submit(fn func()) {
	token := <-gp.tokenChan // 如果没有token 令牌，会阻塞

	go func() {
		fn()
		gp.tokenChan <- token
	}()
}

// Wait 会等待所有任务执行完毕，然后返回
func (gp *GoPool) Wait() {
	for i := 0; i < gp.MaxLimit; i++ {
		<-gp.tokenChan
	}

	close(gp.tokenChan)
}

// Size token 当前大小
func (gp *GoPool) Size() int {
	return len(gp.tokenChan)
}
