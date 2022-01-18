// Copyright (c) 2017-2022 The Elastos Foundation:
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/elastos/Elastos.ELA/account"
	"github.com/elastos/Elastos.ELA/common"
	contract2 "github.com/elastos/Elastos.ELA/core/contract"
	"github.com/elastos/Elastos.ELA/core/contract/program"
	elatypes "github.com/elastos/Elastos.ELA/core/types"
)

func Test_signTransaction(t *testing.T) {
	walletPath := "/Users/github.com/elastos/Elastos.ELA.Rosetta.API/keystore.dat"
	password := "123"
	// from /construction/payloads
	unsignedTransaction := "0902000001ef14bcfc0c0542cf913edb636c016a9c168f5279fa2ce81f58159a54bab1c92900000000000002b037db964a231458d2d6ffd5ea18944c4f90e63d547c5d3b9874df66a4ead0a37a6289cf020000000000000021a796f24c6a36865b83f2546b41553b130b93cfcb00b037db964a231458d2d6ffd5ea18944c4f90e63d547c5d3b9874df66a4ead0a32211b2d40800000000000000218f8fb8b8a5e4648b1ef14a608afeb654d5e205f20000000000"

	txUnsignedBytes, err := hex.DecodeString(unsignedTransaction)
	if err != nil {
		fmt.Printf("decode tx from hexstring err: %s\n", err.Error())
		return
	}

	var txn elatypes.Transaction
	err = txn.DeserializeUnsigned(bytes.NewReader(txUnsignedBytes))
	if err != nil {
		fmt.Printf("deserialize tx err: %s\n", err.Error())
		return
	}

	act, err := account.Open(walletPath, []byte(password))
	if err != nil {
		fmt.Println("open account err:", err)
		return
	}

	redeemScript, err := contract2.CreateStandardRedeemScript(act.GetMainAccount().PublicKey)
	if err != nil {
		log.Printf("pubkey to code err: %s\n", err.Error())
		return
	}
	txn.Programs = []*program.Program{
		&program.Program{
			Code:      redeemScript,
			Parameter: nil,
		},
	}

	tx, err := act.Sign(&txn)
	if err != nil {
		fmt.Println("sign err :", err)
		return
	}

	for i, p := range tx.Programs {
		fmt.Println("signature[", i, "]:", common.BytesToHexString(p.Parameter))
	}
}
