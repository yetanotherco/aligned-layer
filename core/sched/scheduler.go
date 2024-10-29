package sched

import (
	"container/heap"
	"log/slog"
	"time"

	"github.com/yetanotherco/aligned_layer/core/supervisor"
)

type prioQueue []Job

type JobRunner struct {
	incomingJobChannel chan Job
	nextExpiration time.Timer
	queue heap.Interface
}

type Job struct {
	task func() error
	nextRun time.Time
	period time.Duration
	recurrent bool
}

var defaultRunner JobRunner

func init() {
	defaultRunner = NewJobRunner()
}

func NewJobRunner() JobRunner {
	// Initialize to any time, we only want to get an initialized struct
	nextExpiration := time.NewTimer(time.Since(time.Now()))
	// We don't want accidental triggers
	nextExpiration.Stop()

	queue := prioQueue(make([]Job, 0, 1024))
	heap.Init(&queue)

	return JobRunner{
		incomingJobChannel: make(chan Job, 16),
		nextExpiration: *nextExpiration,
		queue: &queue,
	}
}

func (runner *JobRunner) StartRunner() {
	serve := func() {
		runner.serve()
	}
	supervisor.Serve(serve, "job_scheduler")
}

func StartRunner() {
	defaultRunner.StartRunner()
}

func (runner *JobRunner) At(when time.Time, task func() error) {
	job := Job{
		nextRun: when,
		task:    task,
	}
	runner.pushJob(job)
}

func At(when time.Time, task func() error) {
	defaultRunner.At(when, task)
}

func (runner *JobRunner) Every(period time.Duration, task func() error) {
	job := Job{
		nextRun:   time.Now().Add(period),
		task:      task,
		recurrent: true,
		period:    period,
	}
	runner.pushJob(job)
}

func Every(period time.Duration, task func() error) {
	defaultRunner.Every(period, task)
}

func (runner *JobRunner) pushJob(job Job) {
	runner.incomingJobChannel <- job
}

func (runner *JobRunner) runExpiredJobs(now time.Time) {
	runner.nextExpiration.Stop()
	// In Go <1.23 we need to drain the channel for rearming to be safe
	<-runner.nextExpiration.C
	for runner.queue.Len() > 0 {
		job := runner.queue.Pop().(Job)
		if now.Before(job.nextRun) {
			runner.queue.Push(job)
			break
		}
		err := job.task()
		if err != nil {
			slog.Error("periodic task failed", "error", err)
		}
		if job.recurrent {
			job.nextRun = now.Add(job.period)
			runner.queue.Push(job)
		}
	}
	if runner.queue.Len() != 0 {
		// Peek the next job, but DO NOT consume it
		nextJob := runner.queue.Pop().(Job)
		runner.queue.Push(nextJob)
		runner.nextExpiration.Reset(nextJob.nextRun.Sub(now))
	}
}

// Not exported, use the `StartRunner` interface to handle `panic`s
func (runner *JobRunner) serve() {
	// We might be recovering from a crash, so there may be jobs waiting
	// This also resets the timer if there are tasks in the queue
	runner.runExpiredJobs(time.Now())

	for {
		select {
		case newJob := <-runner.incomingJobChannel:
			runner.queue.Push(newJob)
			runner.runExpiredJobs(time.Now())
		case now := <-runner.nextExpiration.C:
			runner.runExpiredJobs(now)
		}
	}
}

// Implementation adapted from: https://pkg.go.dev/container/heap@go1.23.2#example-package-PriorityQueue
func (pq prioQueue) Len() int {
	return len(pq)
}

func (pq prioQueue) Less(i, j int) bool {
	// We want `Pop` to give us the earliest, so we use `After` (`heap.Interface` uses negative prio)
	return pq[i].nextRun.After(pq[j].nextRun)
}

func (pq prioQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *prioQueue) Push(x any) {
	item := x.(Job)
	*pq = append(*pq, item)
}

func (pq *prioQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
