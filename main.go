package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"github.com/hzxiao/neo-thinsdk-go/neo"
	"github.com/hzxiao/neo-thinsdk-go/utils"
	"github.com/hzxiao/neotx/req"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const version = "0.0.2"

var (
	ver     = pflag.BoolP("version", "v", false, "print version")
	help    = pflag.BoolP("help", "h", false, "show usage of neotx command")
	network = pflag.StringP("net", "n", "testnet", "give a network(testnet|mainnet|prinet) of neo")
	node    = pflag.StringP("node", "d", "", "give a node of neo if necessary. it is necessary when network is prinet")
	arg     = pflag.StringP("arg", "a", "arg.json", "a json format file is arg for creating a neo tx")
	send    = pflag.BoolP("send", "s", false, "send rpc request to neo node")
)

func main() {
	pflag.Parse()

	if *ver {
		fmt.Printf("neotx%v\n", version)
		os.Exit(0)
	}
	if *help {
		pflag.Usage()
		os.Exit(0)
	}
	err := req.SetNetwork(*network, *node)
	if err != nil {
		log.Error("[main] set network err: %v", err)
		os.Exit(1)
	}
	err = start(*network, *arg)
	if err != nil {
		log.Error("[main] start err: %v", err)
		os.Exit(1)
	}
}

func start(network string, argFile string) error {
	buf, err := ioutil.ReadFile(argFile)
	if err != nil {
		return err
	}

	arg, err := jsonDecode(buf)
	if err != nil {
		return err
	}

	txParam := &neo.CreateSignParams{}
	txParam.Version = byte(arg.GetInt64("version"))
	txParam.PriKey = arg.GetString("fromPriKey")
	txParam.From = arg.GetString("from")
	txParam.To = arg.GetString("to")
	txParam.AssetId = arg.GetString("assetId")

	factor := uint64(arg.GetInt64("factor"))
	txParam.Value = uint64(arg.GetInt64("value")) * factor

	for _, input := range arg.GetMapArray("input") {
		utxo := neo.Utxo{
			Hash: input.GetString("prevHash"),
			N:    uint16(input.GetInt64("prevIndex")),
		}

		if strings.HasPrefix(utxo.Hash, "0x") {
			utxo.Hash = utxo.Hash[2:]
		}
		utxo.Value, err = req.GetUtxoOutputValue(utxo.Hash, int(utxo.N))
		if err != nil {
			log.Error("[start] get utxo output err: %v", err)
		}
		if utxo.Value == 0 {
			log.Warn("input(%v, %v) value is zero", utxo.Hash, utxo.N)
		}
		utxo.Value = utxo.Value * factor
		txParam.Utxos = append(txParam.Utxos, utxo)
	}

	if len(txParam.Utxos) == 0 {
		hash, _ := neo.GetPublicKeyHashFromAddress(txParam.From)
		txParam.Attrs = append(txParam.Attrs, neo.Attribute{
			Usage: neo.Script,
			Data:  hash,
		})
	}
	now := time.Now().Local().UnixNano() / 1e6
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(now))
	txParam.Attrs = append(txParam.Attrs, neo.Attribute{
		Usage: neo.Remark,
		Data:  b,
	})
	//
	typ := arg.GetString("type")
	var txType byte
	if typ == "InvocationTransaction" {
		invo := arg.GetMap("invocation")
		txParam.Data = neo.InvocationToScript(invo.GetString("contract"), invo.GetString("operation"), invo.GetArray("params"))
		txType = neo.InvocationTransaction
	} else {
		txType = neo.ContractTransaction
	}

	body, raw, err := neo.CreateTx(txType, txParam)
	if err != nil {
		return err
	}

	log.Info("tx body: %v", body)
	log.Info("txid: 0x%v", computeTxid(body))
	log.Info("tx raw:  %v", raw)

	if *send {
		log.Info("send raw tx to network(%v) node(%v)...", network, req.Node)
		err = req.SendRawTransaction(raw)
		if err != nil {
			return fmt.Errorf("send raw tx err: %v", err)
		}
		log.Info("send success")
	}
	return nil
}

func jsonDecode(buf []byte) (goutil.Map, error) {
	var data = goutil.Map{}
	err := json.Unmarshal(buf, &data)
	return data, err
}

func computeTxid(body string) string {
	re := hexlify(string(hash256(unhexlify(body))))

	b, _ := utils.ToBytes(re)

	return utils.ToHexString(utils.BytesReverse(b))
}

func hash256(str string) []byte {
	s := sha256.New()
	s.Write([]byte(str))
	res := s.Sum(nil)
	s.Reset()
	s.Write(res)
	return s.Sum(nil)
}

func hexlify(str string) string {
	return hex.EncodeToString([]byte(str))
}

func unhexlify(str string) string {
	b, _ := hex.DecodeString(str)
	return string(b)
}
