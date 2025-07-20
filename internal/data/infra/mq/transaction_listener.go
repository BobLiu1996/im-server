package mq

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sync"
	"sync/atomic"
	"time"
)

type DemoListener struct {
	localTrans       *sync.Map
	transactionIndex int32
}

func NewDemoListener() *DemoListener {
	return &DemoListener{
		localTrans: new(sync.Map),
	}
}

func (dl *DemoListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	nextIndex := atomic.AddInt32(&dl.transactionIndex, 1)
	fmt.Printf("nextIndex: %v for transactionID: %v\n", nextIndex, msg.TransactionId)
	status := nextIndex % 3
	dl.localTrans.Store(msg.TransactionId, primitive.LocalTransactionState(status+1))

	fmt.Printf("dl")
	time.Sleep(1 * time.Minute)
	return primitive.UnknowState
}

func (dl *DemoListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Printf("%v msg transactionID : %v\n", time.Now(), msg.TransactionId)
	v, existed := dl.localTrans.Load(msg.TransactionId)
	if !existed {
		fmt.Printf("unknown msg: %v, return Commit", msg)
		return primitive.CommitMessageState
	}
	state := v.(primitive.LocalTransactionState)
	switch state {
	case 1:
		fmt.Printf("checkLocalTransaction COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	case 2:
		fmt.Printf("checkLocalTransaction ROLLBACK_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	case 3:
		fmt.Printf("checkLocalTransaction unknown: %v\n", msg)
		return primitive.UnknowState
	default:
		fmt.Printf("checkLocalTransaction default COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	}
}
