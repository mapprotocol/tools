package main

import (
	"context"
	"encoding/json"
	"fmt"
	ethchain "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mapprotocol/atlas/chains/headers/ethereum"
	"github.com/mapprotocol/atlas/cmd/ethclient"
	"github.com/mapprotocol/atlas/core/rawdb"
	"gopkg.in/urfave/cli.v1"
	"log"
	"math/big"
)

func saveMock(ctx *cli.Context) error {
	debugInfo := debugInfo{}
	debugInfo.relayerData = []*relayerInfo{
		&relayerInfo{url: keystore1},
		//{url: keystore2},
		//{url: keystore3},
		//{url: keystore4},
		//{url: keystore5},
		//&relayerInfo{url: keystore6},
	}
	debugInfo.preWork(ctx, true)
	debugInfo.saveMock(ctx) //change this
	return nil
}

func (d *debugInfo) saveMock(ctx *cli.Context) {

	go d.atlasBackend()
	for {
		select {
		case currentEpoch := <-d.notifyCh:
			fmt.Println("CURRENT EPOCH ========>", currentEpoch)
			currentEpoch1 := int(currentEpoch)
			for i := 0; i < len(d.step); i++ {
				if d.step[i] == currentEpoch1 {
					currentEpoch1 = i + 1
					break
				}
			}
			switch currentEpoch1 {
			case 1:
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.queryDebuginfo(QUERY_RELAYERINFO)
				d.queryDebuginfo(BALANCE)
				d.queryDebuginfo(REGISTER_BALANCE)
				d.queryDebuginfo(REWARD)
				d.doSave2(d.ethData2[:250])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.atlasBackendCh <- NEXT_STEP
			case 2:
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.queryDebuginfo(QUERY_RELAYERINFO)
				d.queryDebuginfo(BALANCE)
				d.queryDebuginfo(REGISTER_BALANCE)
				d.queryDebuginfo(REWARD)
				d.doSave(d.ethData[:10])
				d.atlasBackendCh <- NEXT_STEP
			case 3:
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.queryDebuginfo(QUERY_RELAYERINFO)
				d.queryDebuginfo(BALANCE)
				d.queryDebuginfo(REGISTER_BALANCE)
				d.queryDebuginfo(REWARD)
				d.atlasBackendCh <- NEXT_STEP
				return
			default:
				fmt.Println("over")
			}
		}
	}
}

func (d *debugInfo) doSave(chains []ethereum.Header) {
	fmt.Println("=================DO SAVE========================")
	marshal, _ := json.Marshal(chains)
	fmt.Println("len---->", len(marshal))
	conn := d.client
	for k, _ := range d.relayerData {
		fmt.Println("ADDRESS:", d.relayerData[k].from)
		d.relayerData[k].realSave(conn, ChainTypeETH, marshal)
	}
}
func (d *debugInfo) doSave2(chains []types.Header) {
	fmt.Println("=================DO SAVE========================")
	marshal, _ := rlp.EncodeToBytes(chains)
	fmt.Println("len--doSave2-->", len(marshal), "Max:", "131,072")
	conn := d.client
	for k, _ := range d.relayerData {
		fmt.Println("ADDRESS:", d.relayerData[k].from)
		d.relayerData[k].realSave(conn, ChainTypeETH, marshal)
	}
}
func (r *relayerInfo) realSave(conn *ethclient.Client, chainType rawdb.ChainType, marshal []byte) bool {
	//header, err := conn.HeaderByNumber(context.Background(), nil)
	//if err != nil {
	//	log.Fatal(err)
	//	return false
	//}
	input := packInputStore("save", big.NewInt(int64(chainType)), big.NewInt(int64(ChainTypeMAP)), marshal)
	sendContractTransaction(conn, r.from, HeaderStoreAddress, nil, r.priKey, input)

	//input := packInputStore("save", chainType, "MAP", marshal)
	//msg := ethchain.CallMsg{From: r.from, To: &HeaderStoreAddress, Data: input}
	//_, err = conn.CallContract(context.Background(), msg, header.Number)
	//if err != nil {
	//	//log.Fatal("method CallContract error (realSave) :", err)
	//	fmt.Println("save false")
	//	return false
	//}
	//fmt.Println("save success")
	return true
}
func (d *debugInfo) saveByDifferentAccounts(ctx *cli.Context) {
	go d.atlasBackend()
	for {
		select {
		case currentEpoch := <-d.notifyCh:
			fmt.Println("CURRENT EPOCH ========>", currentEpoch)
			currentEpoch1 := int(currentEpoch)
			for i := 0; i < len(d.step); i++ {
				if d.step[i] == currentEpoch1 {
					currentEpoch1 = i + 1
					break
				}
			}
			switch currentEpoch1 {
			case 1:
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				//d.queryDebuginfo(QUERY_RELAYERINFO)
				//d.queryDebuginfo(BALANCE)
				//d.queryDebuginfo(REGISTER_BALANCE)
				//d.query_debugInfo(REWARD)
				d.doSave(d.ethData[:10])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.doSave(d.ethData[:9])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.doSave(d.ethData[:2])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.doSave(d.ethData[:1])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.atlasBackendCh <- NEXT_STEP
			case 2:
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.queryDebuginfo(QUERY_RELAYERINFO)
				d.queryDebuginfo(BALANCE)
				d.queryDebuginfo(REGISTER_BALANCE)
				//d.query_debugInfo(REWARD)
				d.doSave(d.ethData[10:20])
				d.atlasBackendCh <- NEXT_STEP
			case 3:
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.queryDebuginfo(QUERY_RELAYERINFO)
				d.queryDebuginfo(BALANCE)
				d.queryDebuginfo(REGISTER_BALANCE)
				//d.query_debugInfo(REWARD)
				d.doSave(d.ethData[10:20])
				d.atlasBackendCh <- NEXT_STEP
			case 4:
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				d.queryDebuginfo(QUERY_RELAYERINFO)
				d.queryDebuginfo(BALANCE)
				d.queryDebuginfo(REGISTER_BALANCE)
				d.atlasBackendCh <- NEXT_STEP
				return
			default:
				fmt.Println("over")
			}
		}
	}
}

func (d *debugInfo) saveForkBlock(ctx *cli.Context) {
	go d.atlasBackend()
	A, B := getForkBlock()
	//fmt.Println("fmt.Println(A[0].ParentHash)                        ", A[0].ParentHash)
	for {
		select {
		case currentEpoch := <-d.notifyCh:
			fmt.Println("CURRENT EPOCH ========>", currentEpoch)
			currentEpoch1 := int(currentEpoch)
			for i := 0; i < len(d.step); i++ {
				if d.step[i] == currentEpoch1 {
					currentEpoch1 = i + 1
					break
				}
			}
			switch currentEpoch1 {
			case 1:
				//fmt.Println(B[0].Number)
				//a1->a2->a3->a4->...an
				//a1->a2->a3->b4->...bn
				// a1 -> ...a10  want Success
				d.doSave(A[:10])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("1. want Success and is a  Canonical ", A[9].Number, A[9].Hash())

				// b5->   ...b10 want Failed
				d.doSave(B[1:7])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("2.want Failed")

				// b4->   ...b10 want Success but not a  Canonical
				d.doSave(B[0:7])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("3.want Success but not a  Canonical", B[6].Number, B[6].Hash())

				// b4->   ...b11 want Success  is a  Canonical
				d.doSave(B[0:8])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("4.want Success  is a  Canonical: ", B[7].Number, B[7].Hash())

				// a11 -> ...a15 want Success
				d.doSave(A[10:15])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("5.want Success is a Canonical", A[14].Number, A[14].Hash())

				// a11 -> ...a10  want Success but not a  Canonical
				d.doSave(A[:10])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("6.want Success but not a  Canonical", B[9].Number, B[9].Hash())

				// B11 -> ...B15 want Success  is a  Canonical
				d.doSave(B[8:12])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("7.want Success  is a  Canonical", B[11].Number, B[11].Hash(), "  parenthash:  ", B[8].ParentHash)

				// a0 -> ...a16  want Success is a  Canonical
				d.doSave(A[:16])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("8.want Success is a  Canonical", B[15].Number, B[15].Hash())

				// a16 -> ...a20  want Success is a  Canonical
				d.doSave(A[16:20])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("9.want Success is a  Canonical", B[19].Number, B[19].Hash())

				d.atlasBackendCh <- NEXT_STEP
				return
			default:
				fmt.Println("over")
			}
		}
	}
}
func (d *debugInfo) saveForkBlock01(ctx *cli.Context) {
	go d.atlasBackend()
	A, B := getForkBlock()
	//fmt.Println("fmt.Println(A[0].ParentHash)                        ", A[0].ParentHash)
	for {
		select {
		case currentEpoch := <-d.notifyCh:
			fmt.Println("CURRENT EPOCH ========>", currentEpoch)
			currentEpoch1 := int(currentEpoch)
			for i := 0; i < len(d.step); i++ {
				if d.step[i] == currentEpoch1 {
					currentEpoch1 = i + 1
					break
				}
			}
			switch currentEpoch1 {
			case 1:

				// a11 -> ...a15 want Success
				d.doSave(A[0:15])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("1.want Success is a Canonical", A[14].Number, A[14].Hash())

				// b4->   ...b11 want Success  is a  Canonical
				d.doSave(B[0:8])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("2.want Success  not a  Canonical: ", B[7].Number, B[7].Hash())

				// B11 -> ...B15 want Success  is a  Canonical
				d.doSave(B[8:12])
				d.queryDebuginfo(CHAINTYPE_HEIGHT)
				fmt.Println("3.want Success  is a  Canonical", B[11].Number, B[11].Hash(), "  parenthash:  ", B[8].ParentHash)

				d.atlasBackendCh <- NEXT_STEP
				return
			default:
				fmt.Println("over")
			}
		}
	}
}

//  getCurrent type chain number by abi
func getCurrentNumberAbi(conn *ethclient.Client, chainType rawdb.ChainType, from common.Address) (uint64, string) {
	header, err := conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	input := packInputStore(CurNbrAndHash, big.NewInt(int64(chainType)))
	msg := ethchain.CallMsg{From: from, To: &HeaderStoreAddress, Data: input}
	output, err := conn.CallContract(context.Background(), msg, header.Number)
	if err != nil {
		log.Fatal("method CallContract error", err)
	}
	method, _ := abiHeaderStore.Methods[CurNbrAndHash]
	ret, err := method.Outputs.Unpack(output)
	ret1 := ret[0].(*big.Int).Uint64()
	ret2 := common.BytesToHash(ret[1].([]byte))
	return ret1, ret2.String()

}

func packInputStore(abiMethod string, params ...interface{}) []byte {
	input, err := abiHeaderStore.Pack(abiMethod, params...)
	if err != nil {
		log.Fatal(abiMethod, " error ", err)
	}
	return input
}
