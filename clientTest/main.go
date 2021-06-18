package main

import (
	"bufio"
	"container/list"
	"context"
	"fmt"
	"github.com/portto/solana-go-sdk/client"
	"io"
	"os"
	"strconv"
	"strings"
)

const rpcUrl = "http://122.248.192.182:8899"

var arg = map[string]interface{}{
	"encoding":           "json",
	"transactionDetails": "full",
	"rewards":            false,
	"commitment":         "finalized",
}

//两个线程分别从节点获取块数据，对比数据是否一致，验证空块情况
func getAndSaveFile() {

	go func() {
		curSlot := uint64(82939961)

		filePath := "nodelay.txt"
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		defer file.Close()

		client := client.NewClient(rpcUrl)
		for {
			slot, err := client.GetSlot(context.Background())
			if err != nil {
				fmt.Errorf("GetSlot", err)
				break
			}

			for i := curSlot - 1; i < slot; i++ {
				params := []interface{}{i, arg}
				block, err := client.GetConfirmedBlock(context.Background(), params...)
				if err != nil {
					fmt.Errorf("GetConfirmedBlock", err)
					break
				}
				var str string
				if block.Blockhash == "" {
					str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", i, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
				} else {
					str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d\n", i, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
				}
				fmt.Print("no ->" + str)
				file.WriteString(str)
			}
		}
	}()

	curSlot := uint64(82939961)
	filePath := "somedelay.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()

	client := client.NewClient(rpcUrl)
	for {
		slot, err := client.GetSlot(context.Background())
		if err != nil {
			fmt.Errorf("GetSlot", err)
			break
		}
		slot -= 20

		for i := curSlot - 1; i < slot; i++ {
			params := []interface{}{i, arg}
			block, err := client.GetConfirmedBlock(context.Background(), params...)
			if err != nil {
				fmt.Errorf("GetConfirmedBlock", err)
				break
			}
			var str string
			if block.Blockhash == "" {
				str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", i, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
			} else {
				str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d\n", i, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
			}
			fmt.Print("delay ->" + str)
			file.WriteString(str)
		}
	}
}

func getString(str string) string {
	result := strings.Index(str, "EMPTY slot")
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)

		var r = []rune(str)
		length := len(r)
		return string(r[result+11 : length-2])
	}
	return ""
}

func getListFromFile(filePath string) *list.List {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//读原来文件的内容，并且显示在终端
	l := list.New()
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		//fmt.Print(str)

		num := getString(str)
		//fmt.Println(num)
		l.PushBack(num)
	}
	return l
}

func checkEmptyBlock() {
	l := getListFromFile("empty.txt")
	client := client.NewClient(rpcUrl)

	for i := l.Front(); i != nil; i = i.Next() {
		//fmt.Println(i.Value)

		intNum, _ := strconv.Atoi(i.Value.(string))
		height := uint64(intNum)

		params := []interface{}{height, arg}
		block, err := client.GetConfirmedBlock(context.Background(), params...)
		if err != nil {
			fmt.Println("GetConfirmedBlock", err)
			continue
		}
		var str string
		if block.Blockhash == "" {
			str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", height, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
		} else {
			str = fmt.Sprintf("slot:%d ParentSLot:%d PreviousHash:%s size:%d\n", height, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
		}
		fmt.Print(str)
	}
	fmt.Println("Total size: ", l.Len())
}

func checkBlock() {
	client11 := client.NewClient(rpcUrl)
	for {
		slot, err := client11.GetSlot(context.Background())
		if err != nil {
			fmt.Errorf("GetSlot", err)
			break
		}

		params := []interface{}{slot, arg}

		block, err := client11.GetConfirmedBlock(context.Background(), params...)
		if err != nil || block.Blockhash == "" {
			if rpcErr, ok := err.(*client.ErrorResponse); ok {
				//Slot 被跳过打印日志记录，按正常逻辑处理
				//目前发现这两种是可以被忽略的情况,
				if rpcErr.Code == -32007 || rpcErr.Code == -32009 {
					fmt.Println("GetConfirmedBlock", err)
				}
			}
		}

		var str string
		if block.Blockhash == "" {
			str = fmt.Sprintf("slot:%d hash:%s ParentSLot:%d PreviousHash:%s size:%d(EMPTY)\n", slot, block.Blockhash, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
		} else {
			str = fmt.Sprintf("slot:%d hash:%s ParentSLot:%d PreviousHash:%s size:%d\n", slot, block.Blockhash, block.ParentSLot, block.PreviousBlockhash, len(block.Transactions))
		}
		fmt.Print("no ->" + str)
	}
}

func main() {
	checkBlock()
}
