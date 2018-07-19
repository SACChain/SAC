// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"sync"

	"github.com/SACChain/SAC/common"
	"github.com/SACChain/SAC/consensus"
	"github.com/SACChain/SAC/core/types"
	"github.com/SACChain/SAC/log"
)

// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's difficulty requirements.
// 尝试找到一个nonce值能够满足区块难度需求。
// 以上Seal方法体，针对ethash的各种状态进行了校验和流程处理，以及对线程资源的控制
func (ethash *Ethash) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	// If we're running a fake PoW, simply return a 0 nonce immediately
	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModeFullFake {
		header := block.Header()
		header.Nonce, header.MixDigest = types.BlockNonce{}, common.Hash{}
		return block.WithSeal(header), nil
	}
	// If we're running a shared PoW, delegate sealing to it
	// 共享pow的话，则转到它的共享对象执行Seal操作
	if ethash.shared != nil {
		return ethash.shared.Seal(chain, block, stop)
	}
	// Create a runner and the multiple search threads it directs
	// 创建一个runner以及它指挥的多重搜索线程
	abort := make(chan struct{})
	found := make(chan *types.Block)

	ethash.lock.Lock()
	threads := ethash.threads// 挖矿的线程s
	if ethash.rand == nil {// rand为空，则为ethash的字段rand赋值
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))// 获得种子
		if err != nil {
			ethash.lock.Unlock()
			return nil, err
		}
		ethash.rand = rand.New(rand.NewSource(seed.Int64()))// 执行成功，拿到合法种子seed，通过其获得rand对象，赋值。
	}
	ethash.lock.Unlock()
	if threads == 0 {
		threads = runtime.NumCPU()// 挖矿线程编号为0，则通过方法返回当前物理上可用CPU编号
	}
	if threads < 0 {// 非法结果
		threads = 0 // Allows disabling local mining without extra logic around local/remote
		// 置为0，允许在本地或远程没有额外逻辑的情况下，取消本地挖矿操作
	}
	var pend sync.WaitGroup// 创建一个倒计时锁对象
	for i := 0; i < threads; i++ {
		pend.Add(1)
		go func(id int, nonce uint64) {// 核心代码通过闭包多线程技术来执行。
			defer pend.Done()
			ethash.mine(block, id, nonce, abort, found)
		}(i, uint64(ethash.rand.Int63()))//闭包第二个参数表达式uint64(ethash.rand.Int63())通过上面准备好的rand函数随机数结果作为nonce实参传入方法体
	}
	// Wait until sealing is terminated or a nonce is found
	// 直到seal操作被中止或者找到了一个nonce值，否则一直等
	var result *types.Block
	select {
	case <-stop:
		// Outside abort, stop all miner threads
		// 外部意外中止，停止所有挖矿线程
		close(abort)
	case result = <-found:
		// One of the threads found a block, abort all others
		// 其中一个线程挖到正确块，中止其他所有线程
		close(abort)
	case <-ethash.update:
		// Thread count was changed on user request, restart
		// ethash对象发生改变，停止当前所有操作，重启当前方法
		close(abort)
		pend.Wait()
		return ethash.Seal(chain, block, stop)
	}
	// Wait for all miners to terminate and return the block
	// 等待所有矿工停止或者返回一个区块
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block difficulty.
// mine函数是真正的pow矿工，用来搜索一个nonce值，nonce值开始于seed值，
// seed值是能最终产生正确的可匹配可验证的区块难度
func (ethash *Ethash) mine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block) {
	// Extract some data from the header
	// 从区块头中提取出一些数据，放在一个全局变量域中
	var (
		header  = block.Header()
		hash    = header.HashNoNonce().Bytes()
		target  = new(big.Int).Div(maxUint256, header.Difficulty)
		number  = header.Number.Uint64()
		dataset = ethash.dataset(number)
	)
	// Start generating random nonces until we abort or find a good one
	// 开始生成随机nonce值知道我们中止或者成功找到了一个合适的值
	var (
		attempts = int64(0)// 初始化一个尝试次数的变量，下面会利用该变量耍一些花枪
		nonce    = seed// 初始化为seed值，后面每次尝试以后会累加
	)
	logger := log.New("miner", id)
	logger.Trace("Started ethash search for new nonces", "seed", seed)
search:
	for {
		select {
		case <-abort:// 中止命令
			// Mining terminated, update stats and abort
			// 挖矿中止，更新状态，中止当前操作，返回空
			logger.Trace("Ethash nonce search aborted", "attempts", nonce-seed)
			ethash.hashrate.Mark(attempts)
			break search

		default:// 默认执行
			// We don't have to update hash rate on every nonce, so update after after 2^X nonces
			// 我们没必要在每一次尝试nonce值的时候更新hash率，可以在尝试了2的X次方nonce值以后再更新即可
			attempts++
			if (attempts % (1 << 15)) == 0 {// 这里是定的2的15次方
				ethash.hashrate.Mark(attempts)// 满足条件了以后，要更新ethash的hash率字段的状态值
				attempts = 0 // 重置尝试次数
			}
			// Compute the PoW value of this nonce
			// 为这个nonce值计算pow值
			digest, result := hashimotoFull(dataset.dataset, hash, nonce)
			if new(big.Int).SetBytes(result).Cmp(target) <= 0 {
				// Correct nonce found, create a new header with it
				// 找到正确nonce值，创建一个基于它的新的区块头
				header = types.CopyHeader(header)
				header.Nonce = types.EncodeNonce(nonce)// 将输入的整型值转换为一个区块nonce值
				header.MixDigest = common.BytesToHash(digest)// 将字节数组转换为Hash对象

				// Seal and return a block (if still needed)
				// 封装返回一个区块
				select {
				case found <- block.WithSeal(header):
					logger.Trace("Ethash nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
				case <-abort:
					logger.Trace("Ethash nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
				}
				break search
			}
			nonce++// 累加nonce
		}
	}
	// Datasets are unmapped in a finalizer. Ensure that the dataset stays live
	// during sealing so it's not unmapped while being read.
	runtime.KeepAlive(dataset)
}
