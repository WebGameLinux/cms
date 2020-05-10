package rabbitmq

import (
		"sync"
		"time"
)

type WorkerPoolInterface interface {
		Exec(func()) int
		SetWorkerMaxNum(int) WorkerPoolInterface
		Stop(id ...int)
		GetRuns() []int
		Remove(id int)
		Init()
		Name() string
		Destroy()
}

type WorkerPool struct {
		Max                int
		name               string
		RunStatePidTable   []int
		CacheQueue         sync.Map
		CacheNum           int
		CacheMaxNum        int
		Interval           time.Duration
		CheckCacheInterval time.Duration
		RunnerMapper       map[int]*Runner
		sync.RWMutex
}

type Runner struct {
		Ctr  chan bool   // 信号控制通道
		Chan chan func() // 执行器通道
}

func NewWorkPool(max int, name string, cacheMaxNum int, interval time.Duration, checkInterval time.Duration) *WorkerPool {
		var worker = new(WorkerPool)
		worker.Max = max
		worker.Interval = interval
		worker.CheckCacheInterval = checkInterval
		worker.CacheMaxNum = cacheMaxNum
		worker.name = name
		worker.Init()
		return worker
}

func NewRunner() *Runner {
		runner := new(Runner)
		runner.Ctr = make(chan bool, 2)
		runner.Chan = make(chan func())
		return runner
}

func (this *WorkerPool) Exec(task func()) int {
		id := this.GetPid()
		if id < 0 {
				id := this.AddCacheLen()
				if id < 0 {
						return id
				}
				this.CacheQueue.Store(id, task)
				this.RunStatePidTable = append(this.RunStatePidTable, id)
				return id
		}
		this.Add(id)
		this.dispatch(id, task)
		return id
}

func (this *WorkerPool) GetPid() int {
		pid := this.GetFreePid()
		this.Lock()
		defer this.Unlock()
		if pid < 0 && len(this.RunStatePidTable) < this.Max {
				return this.CreateRunner()
		}
		return pid
}

func (this *WorkerPool) GetFreePid() int {
		this.Lock()
		defer this.Unlock()
		var size = len(this.RunnerMapper)
		if size == 0 {
				return 0
		}
		if size+1 < this.Max {
				size++
				return size
		}
		for id, v := range this.RunnerMapper {
				if v == nil {
						continue
				}
				ok := true
				for _, i := range this.RunStatePidTable {
						if id == i {
								ok = false
								break
						}
				}
				if ok {
						return id
				}
		}
		return -1
}

func (this *WorkerPool) CreateRunner() int {
		id := this.CreatePid()
		go this.run(id)
		return id
}

func (this *WorkerPool) CreatePid() int {
		this.Lock()
		defer this.Unlock()
		return len(this.RunStatePidTable) + 1
}

func (this *WorkerPool) dispatch(id int, task func(), times ...int) {
		runner, ok := this.RunnerMapper[id]
		// 移除递归
		if len(times) > 0 && times[0] > 2 {
				return
		}
		if !ok {
				if id >= this.Max {
						return
				}
				go this.run(id)
				defer func(id int, task func(), times int) {
						time.AfterFunc(100*time.Microsecond, func() {
								times++
								this.dispatch(id, task, times)
						})
				}(id, task, times[0])
				return
		}
		// 	唤醒,并重置
		runner.Ctr <- true
		// 投递任务
		runner.Chan <- task
}

func (this *WorkerPool) GetCacheLen() int {
		this.Lock()
		defer this.Unlock()
		return this.CacheNum
}

func (this *WorkerPool) AddCacheLen() int {
		num := this.GetCacheLen()
		this.Lock()
		defer this.Unlock()
		if this.CacheMaxNum <= this.CacheNum {
				return -2
		}
		this.CacheNum = num + 1
		return this.CacheNum
}

func (this *WorkerPool) SubCacheLen() int {
		num := this.GetCacheLen()
		this.Lock()
		defer this.Unlock()
		this.CacheNum = num - 1
		return this.CacheNum
}

func (this *WorkerPool) SetCacheNum(num int) {
		this.Lock()
		defer this.Unlock()
		this.CacheNum = num
}

func (this *WorkerPool) SetWorkerMaxNum(max int) WorkerPoolInterface {
		this.Max = max
		return this
}

func (this *WorkerPool) Stop(id ...int) {
		this.Lock()
		defer this.Unlock()
		if len(id) == 0 {
				for i := 0; i < this.Max; i++ {
						v, ok := this.RunnerMapper[i]
						if !ok {
								continue
						}
						go this.kill(v)
						this.Remove(i)
						delete(this.RunnerMapper, i)
				}
				return
		}
		for _, i := range id {
				v, ok := this.RunnerMapper[i]
				if !ok {
						continue
				}
				go this.kill(v)
				this.Remove(i)
				delete(this.RunnerMapper, i)
		}
}

func (this *WorkerPool) kill(runner *Runner) {
		runner.Ctr <- false
}

func (this *WorkerPool) GetRuns() []int {
		return this.RunStatePidTable
}

func (this *WorkerPool) Remove(id int) {
		this.Lock()
		defer this.Unlock()
		var size = len(this.RunStatePidTable)
		if size == 0 {
				return
		}
		for index := 0; index < size; {
				i := this.RunStatePidTable[index]
				if i != id {
						continue
				}
				if index == 0 {
						if size > 1 {
								this.RunStatePidTable = this.RunStatePidTable[1:]
						} else {
								this.RunStatePidTable = this.RunStatePidTable[0:0]
						}
						continue
				}
				if index > 0 {
						if size <= index+1 {
								this.RunStatePidTable = this.RunStatePidTable[:index]
						} else {
								this.RunStatePidTable = append(this.RunStatePidTable[:index], this.RunStatePidTable[index+1:]...)
						}
						continue
				}
				size = len(this.RunStatePidTable)
				index++
		}
}

func (this *WorkerPool) Init() {
		var table []int
		if this.RunnerMapper != nil && len(this.RunnerMapper) != 0 {
				return
		}
		this.RunnerMapper = make(map[int]*Runner)
		if len(this.RunStatePidTable) != 0 {
				return
		}
		this.RunStatePidTable = table
		this.start()
}

func (this *WorkerPool) start() {
		for i := 0; i < this.GetNum(); i++ {
				go this.run(i)
		}
}

func (this *WorkerPool) GetNum() int {
		return this.Max
}

func (this *WorkerPool) run(id int) {
		times := 0
		runner := NewRunner()
		interval := this.Interval
		if interval == time.Duration(0) {
				interval = 3 * time.Minute
		}
		checkCacheInterval := this.CheckCacheInterval
		if checkCacheInterval == time.Duration(0) {
				checkCacheInterval = 1 * time.Second
		}
		// 每个 pool runner 检查缓存时间分开
		checkCacheInterval = checkCacheInterval + time.Duration(id)*time.Second
		// 检查
		if r, ok := this.RunnerMapper[id]; ok && r != nil {
				runner = r
		} else {
				this.RunnerMapper[id] = runner
		}
		for {
				select {
				case task := <-runner.Chan:
						this.Add(id)
						if task != nil {
								task()
						}
						this.Remove(id)
				case ctr := <-runner.Ctr:
						times = 0
						if !ctr {
								close(runner.Chan)
								goto CtrKill
						}
				case <-time.NewTicker(checkCacheInterval).C:
						task := this.GetCacheTask()
						if task != nil {
								this.Add(id)
								task()
						}
						this.Remove(id)
				case <-time.NewTicker(interval).C:
						times++
				}
				if times > this.CacheMaxNum*3 {
						this.Stop(id)
				}
		}
CtrKill:
}

func (this *WorkerPool) Add(id int) {
		this.Lock()
		defer this.Unlock()
		for _, i := range this.RunStatePidTable {
				if i == id {
						return
				}
		}
		this.RunStatePidTable = append(this.RunStatePidTable, id)
}

func (this *WorkerPool) GetCacheTask() func() {
		var (
				ok   bool
				task func()
		)
		this.CacheQueue.Range(func(key, value interface{}) bool {
				if key == nil {
						return true
				}
				if value == nil {
						return true
				}
				if task, ok = value.(func()); ok {
						this.CacheQueue.Delete(key)
						return false
				}
				return true
		})
		if task != nil {
				this.SubCacheLen()
		}
		return task
}

func (this *WorkerPool) Name() string {
		return this.name
}

func (this *WorkerPool) Destroy() {
		this.CacheNum = 0
		this.CacheQueue.Range(func(key, value interface{}) bool {
				this.CacheQueue.Delete(key)
				this.SubCacheLen()
				return true
		})
		for _, id := range this.RunStatePidTable {
				this.Stop(id)
		}
		this.RunStatePidTable = this.RunStatePidTable[0:0]
}
