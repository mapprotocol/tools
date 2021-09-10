package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	ethchain "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mapprotocol/atlas/chains/headers/ethereum"
	"github.com/mapprotocol/atlas/core/rawdb"
	"github.com/mapprotocol/atlas/core/vm"
	"github.com/mapprotocol/atlas/params"
	params2 "github.com/mapprotocol/atlas/params"
	"github.com/mapprotocol/tools/ethclient"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"time"
)

var (
	epochHeight                       = params.NewEpochLength
	keystore1                         = "D:/work/atlas_master/atlas/data/keystore/UTC--2021-07-19T02-04-57.993791200Z--df945e6ffd840ed5787d367708307bd1fa3d40f4"
	keystore2                         = "D:/BaiduNetdiskDownload/test015/atlas-fork/atlas/data555/keystore/UTC--2021-07-09T06-26-32.960000300Z--78c5285c42572677d3f9dcc27b9ac7b1ff49843c"
	keystore3                         = "D:/BaiduNetdiskDownload/test015/atlas-fork/atlas/data555/keystore/UTC--2021-07-11T06-35-36.635750800Z--70bf8d9de50713101992649a4f0d7fa505ebb334"
	keystore4                         = "D:/BaiduNetdiskDownload/test015/atlas-fork/atlas/data555/keystore/UTC--2021-07-19T11-51-51.704095400Z--4e0449459f73341f8e9339cb9e49dae3115ec80f"
	keystore5                         = "D:/BaiduNetdiskDownload/test015/atlas-fork/atlas/data555/keystore/UTC--2021-07-21T10-26-12.236878500Z--8becddb5fbe6f3d6b08450e2d33e48e63d6c4b29"
	keystore6                         = "D:/BaiduNetdiskDownload/test015/atlas-fork/atlas/data555/keystore/UTC--2021-08-08T07-06-17.823389800Z--4c179dd018ab2852bb3b76f4e3c26de797997601"
	password                          = "111111"
	abiRelayer, _                     = abi.JSON(strings.NewReader(params2.RelayerABIJSON))
	abiHeaderStore, _                 = abi.JSON(strings.NewReader(params2.HeaderStoreABIJSON))
	RelayerAddress     common.Address = params2.RelayerAddress
	HeaderStoreAddress common.Address = params2.HeaderStoreAddress
)

const (
	BALANCE           = "balance"
	REGISTER_BALANCE  = "registerBalance"
	QUERY_RELAYERINFO = "relayerInfo"
	REWARD            = "reward"
	CHAINTYPE_HEIGHT  = "chainTypeHeight"
	NEXT_STEP         = "next step"

	AtlasRPCListenAddr = "localhost" //map 119.8.165.158 //pist 3.35.104.123
	AtlasRPCPortFlag   = 8082

	EthRPCListenAddr = "119.8.165.158"
	EthRPCPortFlag   = 8545
	EthUrl           = "https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"

	ChainTypeETH = 3
	ChainTypeMAP = 211

	// method name
	CurNbrAndHash = vm.CurNbrAndHash
)

type step []int // epoch height
type debugInfo struct {
	atlasBackendCh chan string
	notifyCh       chan uint64
	step           step
	ethData        []ethereum.Header
	ethData2       []types.Header
	client         *ethclient.Client
	relayerData    []*relayerInfo
}
type relayerInfo struct {
	url           string
	from          common.Address
	preBalance    *big.Float
	nowBalance    *big.Float
	registerValue int64
	priKey        *ecdsa.PrivateKey
	fee           uint64
}

func (r *relayerInfo) swapBalance() {
	f, _ := (r.nowBalance).Float64()
	r.preBalance = big.NewFloat(f)
}

func (r *relayerInfo) changeRegisterValue(value int64) {
	r.registerValue = value
}
func (d *debugInfo) changeAllRegisterValue(value int64) {
	for k, _ := range d.relayerData {
		d.relayerData[k].registerValue = value
	}
}

func (d *debugInfo) preWork(ctx *cli.Context, isRegister bool) {
	conn := getConn11(ctx)
	d.atlasBackendCh = make(chan string)
	d.notifyCh = make(chan uint64)
	d.client = conn

	//d.ethData = getEthChains()
	d.ethData2 = getEthChains2()
	number, err := conn.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("get BlockNumber err ", err)
	}
	currentEpoch := number / epochHeight
	d.step = []int{int(currentEpoch + 1), int(currentEpoch + 2), int(currentEpoch + 3)}
	for k, _ := range d.relayerData {
		Ele := d.relayerData[k]
		priKey, from := loadprivateCommon(Ele.url)
		var acc common.Address
		acc.SetBytes(from.Bytes())
		Ele.registerValue = registerValue
		Ele.from = acc
		Ele.priKey = priKey
		Ele.fee = uint64(0)
		bb := getBalance(conn, Ele.from)
		Ele.preBalance = bb
		Ele.nowBalance = bb
		//if isRegister {
		//	register11(ctx, d.client, *d.relayerData[k])
		//}
	}

}
func (d *debugInfo) queryDebuginfo(ss string) {
	conn := d.client
	switch ss {
	case BALANCE:
		for k, _ := range d.relayerData {
			fmt.Println("ADDRESS:", d.relayerData[k].from, " OLD BALANCE :", d.relayerData[k].preBalance, " NOW BALANCE :", getBalance(conn, d.relayerData[k].from))
		}
	case REGISTER_BALANCE:
		for k, _ := range d.relayerData {
			registered, unregistering, unregistered := getRegisterBalance(conn, d.relayerData[k].from)
			fmt.Println("ADDRESS:", d.relayerData[k].from,
				" NOW registerValue BALANCE :", registered, " registerING BALANCE :", unregistering, "registerED BALANCE :", unregistered)
		}
	case QUERY_RELAYERINFO:
		for k, _ := range d.relayerData {
			bool1, bool2, relayerEpoch, _ := queryRegisterInfo(conn, d.relayerData[k].from)
			fmt.Println("ADDRESS:", d.relayerData[k].from, "ISREGISTER:", bool1, " ISRELAYER :", bool2, " RELAYER_EPOCH :", relayerEpoch)
		}
	case REWARD:

	case CHAINTYPE_HEIGHT:
		for k, _ := range d.relayerData {
			currentTypeHeight, hash := getCurrentNumberAbi(conn, ChainTypeETH, d.relayerData[k].from)
			fmt.Println("ADDRESS:", d.relayerData[k].from, " TYPE HEIGHT:", currentTypeHeight, "  HASH:   ", hash)
		}
	}

}
func (d *debugInfo) atlasBackend() {
	canNext := "YES"
	count := 0
	conn := d.client
	var target uint64 // 1 2 3...
	target = uint64(d.step[count]) - 1
	go func() {
		for {
			select {
			case <-d.atlasBackendCh:
				count++
				if count < len(d.step) {
					target = uint64(d.step[count]) - 1
				}

				canNext = "YES"
			}
		}
	}()

	for {
		number, err := conn.BlockNumber(context.Background())
		if err != nil {
			log.Fatal("get BlockNumber err ", err)
		}
		if canNext != "NO" {
			temp := int(number) - int(target*epochHeight)
			if temp > 0 {
				if int((target+1)*epochHeight) < int(number) {
					log.Fatal("Conditions can never be met")
				}
				d.notifyCh <- target + 1
				canNext = "NO"
				if count+1 == len(d.step) {
					return
				}
			}
		}
		time.Sleep(time.Second)
	}
}
func getEthChains() []ethereum.Header {
	Db, err := rawdb.NewLevelDBDatabase("xxxxxx", 128, 1024, "", false)
	if err != nil {
		log.Fatal(err)
	}
	var key []byte
	key = []byte("ETH_INFO")
	var c []ethereum.Header
	jsonbyte, err := Db.Get(key)
	json.Unmarshal(jsonbyte, &c)
	if len(c) == 1000 {
		return c
	}
	Ethconn, _ := dialEthConn()
	Headers := getChainsCommon(Ethconn)

	rlp, err := json.Marshal(Headers)
	if err != nil {
		log.Fatal("Failed to Marshal block body", "err", err)
	}

	if err := Db.Put(key, rlp); err != nil {
		log.Fatal("Failed to store block body", "err", err)
	}
	return Headers
}
func getEthChains2() []types.Header {
	//Db, err := rawdb.NewLevelDBDatabase("xxxxxx", 128, 1024, "", false)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//var key []byte
	//key = []byte("ETH_INFO2")
	//var c []types.Header
	//jsonbyte, err := Db.Get(key)
	//json.Unmarshal(jsonbyte, &c)
	//if len(c) == 1000 {
	//	return c
	//}
	Ethconn, _ := dialEthConn()
	Headers := getChainsCommon_ethHeaders(Ethconn)

	//rlp, err := json.Marshal(Headers)
	//if err != nil {
	//	log.Fatal("Failed to Marshal block body", "err", err)
	//}

	//if err := Db.Put(key, rlp); err != nil {
	//	log.Fatal("Failed to store block body", "err", err)
	//}
	return Headers
}
func getChainsCommon(conn *ethclient.Client) []ethereum.Header {
	startNum := 1
	endNum := 1000
	Headers := make([]ethereum.Header, 1000)
	HeaderBytes := make([]bytes.Buffer, 1000)
	for i := startNum; i <= endNum; i++ {
		Header, err := conn.HeaderByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			log.Fatal(err)
		}
		convertChain(&Headers[i-1], &HeaderBytes[i-1], Header)
	}
	return Headers
}

func getChainsCommon_ethHeaders(conn *ethclient.Client) []types.Header {
	startNum := 10983101
	endNum := 10983122
	Headers := make([]types.Header, 22)
	j := 0
	for i := startNum; i <= endNum; i++ {
		Header, err := conn.HeaderByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			log.Fatal(err)
		}
		Headers[j] = *Header
		j++
	}
	return Headers
}
func ethToWei(registerValue int64) *big.Int {
	baseUnit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	value := new(big.Int).Mul(big.NewInt(registerValue), baseUnit)
	return value
}

func weiToEth(value *big.Int) uint64 {
	baseUnit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	valueT := new(big.Int).Div(value, baseUnit).Uint64()
	return valueT
}
func printChangeBalance(old, new big.Float) {
	f, _ := old.Float64()
	old1 := big.NewFloat(f)
	f2, _ := new.Float64()
	new1 := big.NewFloat(f2)
	f3, _ := old1.Float64()
	c := big.NewFloat(f3)
	fmt.Printf("old balance:%v  new balance %v  change %v\n",
		old1, new1, c.Abs(c.Sub(c, new1)))
}
func getBalance(conn *ethclient.Client, address common.Address) *big.Float {
	//balance, err := conn.BalanceAt(context.Background(), address, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//balance2 := new(big.Float)
	//balance2.SetString(balance.String())
	//Value := new(big.Float).Quo(balance2, big.NewFloat(math.Pow10(18)))
	return big.NewFloat(0)
}
func getRegisterBalance(conn *ethclient.Client, from common.Address) (uint64, uint64, uint64) {
	header, err := conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	input := packInput("getRelayerBalance", from)
	msg := ethchain.CallMsg{From: from, To: &RelayerAddress, Data: input}
	output, err := conn.CallContract(context.Background(), msg, header.Number)
	if err != nil {
		log.Fatal("method CallContract error", err)
	}
	method, _ := abiRelayer.Methods["getRelayerBalance"]
	ret, err := method.Outputs.Unpack(output)
	if len(ret) != 0 {
		args := struct {
			registered    *big.Int
			unregistering *big.Int
			unregistered  *big.Int
		}{
			ret[0].(*big.Int),
			ret[1].(*big.Int),
			ret[2].(*big.Int),
		}
		return weiToEth(args.registered), weiToEth(args.unregistering), weiToEth(args.unregistered)
	}
	log.Fatal("Contract query failed result len == 0")
	return 0, 0, 0
}
func dialEthConn() (*ethclient.Client, string) {
	//ip = EthRPCListenAddr //utils.RPCListenAddrFlag.Name)
	//port = EthRPCPortFlag //utils.RPCPortFlag.Name)
	//url := fmt.Sprintf("http://%s", fmt.Sprintf("%s:%d", ip, port))
	url := EthUrl
	conn, err := ethclient.Dial(EthUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the AtlasChain client: %v", err)
	}
	return conn, url
}
func register11(ctx *cli.Context, conn *ethclient.Client, info relayerInfo) {
	value := ethToWei(info.registerValue)
	if info.registerValue < RegisterAmount {
		log.Fatal("Amount must bigger than ", RegisterAmount)
	}
	fee := ctx.GlobalUint64(FeeFlag.Name)
	checkFee(new(big.Int).SetUint64(fee))
	input := packInput("register", value)
	sendContractTransaction(conn, info.from, RelayerAddress, nil, info.priKey, input)
}
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func convertChain(header *ethereum.Header, headerbyte *bytes.Buffer, e *types.Header) (*ethereum.Header, *bytes.Buffer) {
	if header == nil || e == nil {
		fmt.Println("header:", header, "e:", e)
		return header, headerbyte
	}
	header.ParentHash = e.ParentHash
	header.UncleHash = e.UncleHash
	header.Coinbase = e.Coinbase
	header.Root = e.Root
	header.TxHash = e.TxHash
	header.ReceiptHash = e.ReceiptHash
	header.GasLimit = e.GasLimit
	header.GasUsed = e.GasUsed
	header.Time = e.Time
	header.MixDigest = e.MixDigest
	header.Nonce = types.EncodeNonce(e.Nonce.Uint64())
	header.Bloom.SetBytes(e.Bloom.Bytes())
	if header.BaseFee = new(big.Int); e.BaseFee != nil {
		header.BaseFee.Set(e.BaseFee)
	}
	if header.Difficulty = new(big.Int); e.Difficulty != nil {
		header.Difficulty.Set(e.Difficulty)
	}
	if header.Number = new(big.Int); e.Number != nil {
		header.Number.Set(e.Number)
	}
	if len(e.Extra) > 0 {
		header.Extra = make([]byte, len(e.Extra))
		copy(header.Extra, e.Extra)
	}
	binary.Write(headerbyte, binary.BigEndian, header)
	return header, headerbyte
}
func registerCommon(conn *ethclient.Client, keystore string) (*big.Float, common.Address) {
	fee := uint64(0)
	value := ethToWei(100000)
	priKey, from := loadprivateCommon(keystore)

	pkey, pk, _ := getPubKey(priKey)
	aBalance := getBalance(conn, from)
	fmt.Printf("Fee: %v \nPub key:%v\nvalue:%v\n \n", fee, pkey, value)
	input := packInput("register", pk, new(big.Int).SetUint64(fee), value)
	sendContractTransaction(conn, from, RelayerAddress, nil, priKey, input)
	return aBalance, from
}

func loadprivateCommon(keyfile string) (*ecdsa.PrivateKey, common.Address) {
	keyjson, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read the keyfile at '%s': %v", keyfile, err))
	}
	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		log.Fatal(fmt.Errorf("error decrypting key: %v", err))
	}
	priKey1 := key.PrivateKey
	return priKey1, crypto.PubkeyToAddress(priKey1.PublicKey)
}

//0xb2be7149dca8d266b0a6f76cc0308844865ffb4320419854720870074c8d1c90
func getBase64(num int64) {
	conn, _ := dialEthConn()
	h, _ := conn.HeaderByNumber(context.Background(), big.NewInt(num))
	encodedStr := base64.StdEncoding.EncodeToString(h.Extra)
	fmt.Println(encodedStr)
}

func getJson(num int64) {
	conn, _ := dialEthConn()
	h, _ := conn.HeaderByNumber(context.Background(), big.NewInt(num))
	Header := &ethereum.Header{}
	Buffer := &bytes.Buffer{}
	h2, _ := convertChain(Header, Buffer, h)
	bs, err := json.Marshal(h2)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
	//encodedStr := base64.StdEncoding.EncodeToString(h.Extra)
	//fmt.Println(encodedStr)
}
