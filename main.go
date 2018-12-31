package main

import (
	"encoding/json"
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"github.com/hzxiao/neo-thinsdk-go/neo"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
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

	raw, err := neo.CreateTx(txType, txParam)
	if err != nil {
		return err
	}
	println(raw)

	return nil
}

func jsonDecode(buf []byte) (goutil.Map, error) {
	var data = goutil.Map{}
	err := json.Unmarshal(buf, &data)
	return data, err
}
