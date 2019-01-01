package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"github.com/hzxiao/neo-thinsdk-go/neo"
	"github.com/hzxiao/neo-thinsdk-go/utils"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"strings"
)

const version = "0.0.1"

var (
	ver     = pflag.BoolP("version", "v", false, "print version")
	help    = pflag.BoolP("help", "h", false, "show usage of neotx command")
	network = pflag.StringP("net", "n", "testnet", "give a network(testnet|mainnet) of neo")
	arg     = pflag.StringP("arg", "a", "arg.json", "a json format file is arg for creating a neo tx")
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

	err := start(*network, *arg)
	if err != nil {
		log.Error("[main] start err: %v", err)
		os.Exit(1)
	}
}

func start(network string, argFile string) error {
	if network != "testnet" && network != "mainnet" {
		return fmt.Errorf("unknown network of neo")
	}

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
	txParam.Value = uint64(arg.GetInt64("value"))

	for _, input := range arg.GetMapArray("input") {
		utxo := neo.Utxo{
			Hash:  input.GetString("hash"),
			Value: uint64(arg.GetInt64("value")),
			N:     uint16(arg.GetInt64("n")),
		}

		if strings.HasPrefix(utxo.Hash, "0x") {
			utxo.Hash = utxo.Hash[2:]
		}
		txParam.Utxos = append(txParam.Utxos, utxo)
	}

	//
	typ := arg.GetString("type")
	var txType byte
	if typ == "InvocationTransaction" {
		invo := arg.GetMap("invocation")
		txParam.Data = neo.InvocationToScript(invo.GetString("contract"), arg.GetString("operation"), arg.GetArray("params"))

		txType = neo.InvocationTransaction

	} else {
		txType = neo.ContractTransaction
	}

	body, raw, err := neo.CreateTx(txType, txParam)
	if err != nil {
		return err
	}

	fmt.Printf("txid: 0x%v\n", computeTxid(body))
	fmt.Printf("tx raw: %v\n", raw)
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
