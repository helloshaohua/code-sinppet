package main

import "log"

type OnceCron struct {
	tasks  []*Task       // 任务的队列
	add    chan *Task    // 当遇到新任务的时候
	remove chan string   // 当遇到删除任务的时候
	stop   chan struct{} // 当遇到停止信号的时候
	Logger *log.Logger   // 日志
}

type Job interface {
	Run()
}

type Task struct {
	Job     Job    // 要执行的任务
	Uuid    string // 任务标识，删除时用
	RunTime int64  // 执行时间
	Spacing int64  // 间隔时间
	EndTime int64  // 结束时间
	Number  int
}

func main() {

}
