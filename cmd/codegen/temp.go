package main

import "time"

type CommentSubject struct {
	Id         int64     `json:"id"`
	ObjId      int64     `json:"objId"`
	ObjType    int8      `json:"objType"`
	MemberId   uint64    `json:"memberId"`
	Count      int32     `json:"count"`
	RootCount  int32     `json:"rootCount"`
	AllCount   int32     `json:"allCount"`
	State      int8      `json:"state"`
	Attrs      int32     `json:"attrs"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}
