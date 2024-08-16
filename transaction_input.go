package main

import "bytes"

// TXInput 交易输入, 输入引用来自一个先前交易的输出
type TXInput struct {
	Txid      []byte // 交易的ID
	Vout      int    // 交易中输出的索引
	Signature []byte // 提供输出中脚本SCriptPubKey使用的数据，如果数据正确，可以解锁输出，并且其值可以用于生成新的输出,如果数据不正确，则无法在输入中引用输出。这是保证用户不能花费属于他人币的机制。
	PubKey    []byte
}

// UseKey 检查输入是否使用特定秘钥来解锁输出
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
