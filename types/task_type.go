package types

import "context"

// 任务信息
type Task struct {
	ID      string
	Name    string
	Status  string
	Message chan string
	Cancel  context.CancelFunc
	Ctx     context.Context
	Type    int
	Data    *TaskData
}

// 添加任务传递的参数
type TaskParams struct {
	ID              string
	Name            string
	Status          string
	Type            int
	AccountClassify string
	DataClassify    []int
	Keyword         string
	TwitterId       string
	Compensate      int
}

// 任务携带的数据
type TaskData struct {
	AccountClassify string
	DataClassify    []int
	Keyword         string
	TwitterId       string
	Compensate      int
}
