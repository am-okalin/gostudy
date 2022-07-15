package singleton

import "sync"

//Message 单例存储的结构体，如数据库句柄，config等
type Message struct {
	Count int
	Msg   string
}

//MessagePool 单例实现的结构体
type MessagePool struct {
	Pool *sync.Pool
}

// 饿汉
var hunger = &MessagePool{Pool: &sync.Pool{New: func() interface{} {
	return &Message{Count: 0}
}}}

//lazy 懒汉 延迟加载(lazy loading) 线程不安全
var lazy *MessagePool

func instance1() *MessagePool {
	return hunger
}

//instance2 线程不安全
func instance2() *MessagePool {
	if lazy == nil {
		lazy = &MessagePool{Pool: &sync.Pool{New: func() interface{} {
			return &Message{Count: 0}
		}}}
	}
	return lazy
}

var lock = sync.Mutex{}

// instance3 对共享资源加锁保证线程安全
func instance3() *MessagePool {
	lock.Lock()
	defer lock.Unlock()
	if lazy == nil {
		lazy = &MessagePool{Pool: &sync.Pool{New: func() interface{} {
			return &Message{Count: 0}
		}}}
	}
	return lazy
}

var once = &sync.Once{}

// instance4 仅执行一次，保证线程安全
func instance4() *MessagePool {
	once.Do(func() {
		lazy = &MessagePool{Pool: &sync.Pool{New: func() interface{} {
			return &Message{Count: 0}
		}}}
	})

	return lazy
}
