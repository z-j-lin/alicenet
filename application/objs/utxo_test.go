package objs

import (
	"bytes"
	"testing"

	"github.com/alicenet/alicenet/application/objs/uint256"
	"github.com/alicenet/alicenet/constants"
	"github.com/alicenet/alicenet/crypto"
)

func TestUTXOCreateValueStore(t *testing.T) {
	chainID := uint32(0)
	value, err := new(uint256.Uint256).FromUint64(65537)
	if err != nil {
		t.Fatal(err)
	}
	fee, err := new(uint256.Uint256).FromUint64(0)
	if err != nil {
		t.Fatal(err)
	}
	acct := make([]byte, constants.OwnerLen)
	curveSpec := constants.CurveSecp256k1
	txHash := make([]byte, constants.HashLen)
	utxo := &TXOut{}
	err = utxo.CreateValueStore(chainID, value, fee, acct, curveSpec, txHash)
	if err == nil {
		t.Fatal("Should have raised error")
	}
	chainID = 1
	err = utxo.CreateValueStore(chainID, value, fee, acct, curveSpec, txHash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTXOCreateValueStoreFromDeposit(t *testing.T) {
	chainID := uint32(0)
	value, err := new(uint256.Uint256).FromUint64(65537)
	if err != nil {
		t.Fatal(err)
	}
	acct := make([]byte, constants.OwnerLen)
	nonce := make([]byte, constants.HashLen)
	utxo := &TXOut{}
	err = utxo.CreateValueStoreFromDeposit(chainID, value, acct, nonce)
	if err == nil {
		t.Fatal("Should have raised error")
	}
	chainID = 1
	err = utxo.CreateValueStoreFromDeposit(chainID, value, acct, nonce)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTXODataStoreGood(t *testing.T) {
	cid := uint32(2)
	idxPre := []byte("Index")
	idx := crypto.Hasher(idxPre)
	txoid := uint32(17)
	iat := uint32(1)
	rawdata := make([]byte, 1)
	dep, err := new(uint256.Uint256).FromUint64(uint64((len(rawdata) + constants.BaseDatasizeConst) * 3))
	if err != nil {
		t.Fatal(err)
	}

	ownerSigner := &crypto.Secp256k1Signer{}
	if err := ownerSigner.SetPrivk(crypto.Hasher([]byte("a"))); err != nil {
		t.Fatal(err)
	}
	ownerPubk, err := ownerSigner.Pubkey()
	if err != nil {
		t.Fatal(err)
	}

	ownerAcct := crypto.GetAccount(ownerPubk)
	owner := &DataStoreOwner{}
	owner.New(ownerAcct, constants.CurveSecp256k1)

	dsp := &DSPreImage{
		ChainID:  cid,
		Index:    idx,
		IssuedAt: iat,
		Deposit:  dep,
		RawData:  rawdata,
		TXOutIdx: txoid,
		Owner:    owner,
		Fee:      new(uint256.Uint256).SetZero(),
	}
	txHash := make([]byte, constants.HashLen)
	dsl := &DSLinker{
		DSPreImage: dsp,
		TxHash:     txHash,
	}
	ds := &DataStore{
		DSLinker: dsl,
	}
	err = ds.PreSign(ownerSigner)
	if err != nil {
		t.Fatal(err)
	}
	utxo := &TXOut{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}

	if !utxo.HasDataStore() {
		t.Fatal("Should have DataStore!")
	}
	if utxo.HasValueStore() {
		t.Fatal("Should not have ValueStore!")
	}

	dsCopy, err := utxo.DataStore()
	if err != nil {
		t.Fatal(err)
	}
	dsEqual(t, ds, dsCopy)

	_, err = utxo.ValueStore()
	if err == nil {
		t.Fatal("Should raise error for no ValueStore!")
	}

	utxo2 := &TXOut{}
	utxoBytes, err := utxo.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	err = utxo2.UnmarshalBinary(utxoBytes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTXOValueStoreGood(t *testing.T) {
	cid := uint32(2)
	val, err := new(uint256.Uint256).FromUint64(65537)
	if err != nil {
		t.Fatal(err)
	}
	txoid := uint32(17)

	ownerSigner := &crypto.Secp256k1Signer{}
	if err := ownerSigner.SetPrivk(crypto.Hasher([]byte("a"))); err != nil {
		t.Fatal(err)
	}
	ownerPubk, err := ownerSigner.Pubkey()
	if err != nil {
		t.Fatal(err)
	}

	ownerAcct := crypto.GetAccount(ownerPubk)
	owner := &ValueStoreOwner{}
	owner.New(ownerAcct, constants.CurveSecp256k1)

	vsp := &VSPreImage{
		ChainID:  cid,
		Value:    val,
		TXOutIdx: txoid,
		Owner:    owner,
		Fee:      new(uint256.Uint256).SetZero(),
	}
	txHash := make([]byte, constants.HashLen)
	vs := &ValueStore{
		VSPreImage: vsp,
		TxHash:     txHash,
	}

	utxo := &TXOut{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}

	if utxo.HasDataStore() {
		t.Fatal("Should not have DataStore!")
	}
	if !utxo.HasValueStore() {
		t.Fatal("Should have ValueStore!")
	}

	_, err = utxo.DataStore()
	if err == nil {
		t.Fatal("Should raise error for no DataStore!")
	}

	vsCopy, err := utxo.ValueStore()
	if err != nil {
		t.Fatal(err)
	}
	vsEqual(t, vs, vsCopy)

	utxo2 := &TXOut{}
	utxoBytes, err := utxo.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	err = utxo2.UnmarshalBinary(utxoBytes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTXOMarshalBinary(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.MarshalBinary()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}
	_, err = utxo.MarshalBinary()
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.MarshalBinary()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.MarshalBinary()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOUnmarshalBinary(t *testing.T) {
	utxo := &TXOut{}
	data := make([]byte, 0)
	err := utxo.UnmarshalBinary(data)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestUTXOPreHash(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.PreHash()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	_, err = utxo.PreHash()
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.PreHash()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.PreHash()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOUTXOID(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.UTXOID()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.UTXOID()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.UTXOID()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOChainID(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.ChainID()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.ChainID()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.ChainID()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOTxOutIdx(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.TxOutIdx()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.TxOutIdx()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.TxOutIdx()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOSetTxOutIdx(t *testing.T) {
	idx := uint32(0)
	utxo := &TXOut{}
	err := utxo.SetTxOutIdx(idx)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.SetTxOutIdx(idx)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.SetTxOutIdx(idx)
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOTxHash(t *testing.T) {
	txHashTrue := make([]byte, constants.HashLen)
	utxo := &TXOut{}
	_, err := utxo.TxHash()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.TxHash()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}
	ds.DSLinker = &DSLinker{}
	ds.DSLinker.TxHash = txHashTrue
	txHash, err := utxo.TxHash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(txHash, txHashTrue) {
		t.Fatal("txHash does not match (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.TxHash()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
	vs.VSPreImage = &VSPreImage{}
	vs.TxHash = txHashTrue
	txHash, err = utxo.TxHash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(txHash, txHashTrue) {
		t.Fatal("txHash does not match (4)")
	}
}

func TestUTXOSetTxHash(t *testing.T) {
	txHash := make([]byte, 0)
	utxo := &TXOut{}
	err := utxo.SetTxHash(txHash)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	txHash = make([]byte, constants.HashLen)
	err = utxo.SetTxHash(txHash)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.SetTxHash(txHash)
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.SetTxHash(txHash)
	if err == nil {
		t.Fatal("Should have raised error (5)")
	}
}

func TestUTXOIsExpired(t *testing.T) {
	currentHeight := uint32(1)
	utxo := &TXOut{}
	_, err := utxo.IsExpired(currentHeight)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.IsExpired(currentHeight)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	expired, err := utxo.IsExpired(currentHeight)
	if err != nil {
		t.Fatal(err)
	}
	if expired {
		t.Fatal("ValueStore should not be expired")
	}
}

func TestUTXORemainingValue(t *testing.T) {
	utxo := &TXOut{}
	currentHeight := uint32(0)
	_, err := utxo.RemainingValue(currentHeight)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.RemainingValue(currentHeight)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.RemainingValue(currentHeight)
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOMakeTxIn(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.MakeTxIn()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.MakeTxIn()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.MakeTxIn()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOValue(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.Value()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.Value()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.Value()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOValuePlusFee(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.ValuePlusFee()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.ValuePlusFee()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.ValuePlusFee()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOValidatePreSignature(t *testing.T) {
	utxo := &TXOut{}
	err := utxo.ValidatePreSignature()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.ValidatePreSignature()
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.ValidatePreSignature()
	if err != nil {
		t.Fatal("Should pass (2)")
	}
}

func TestUTXOValidateSignature(t *testing.T) {
	utxo := &TXOut{}
	currentHeight := uint32(0)
	txIn := &TXIn{}
	err := utxo.ValidateSignature(currentHeight, txIn)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.ValidateSignature(currentHeight, txIn)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.ValidateSignature(currentHeight, txIn)
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestUTXOMustBeMinedBeforeHeight(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.MustBeMinedBeforeHeight()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	iat := uint32(1)
	heightTrue := iat*constants.EpochLength - 1

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.MustBeMinedBeforeHeight()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}
	ds.DSLinker = &DSLinker{}
	ds.DSLinker.DSPreImage = &DSPreImage{}
	ds.DSLinker.DSPreImage.IssuedAt = iat
	height, err := utxo.MustBeMinedBeforeHeight()
	if err != nil {
		t.Fatal(err)
	}
	if height != heightTrue {
		t.Fatal("Incorrect MinedBefore (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.MustBeMinedBeforeHeight()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTXOAccount(t *testing.T) {
	acct := make([]byte, constants.OwnerLen)
	curveSpec := constants.CurveSecp256k1
	o := &Owner{}
	err := o.New(acct, curveSpec)
	if err != nil {
		t.Fatal(err)
	}
	altOwner := &Owner{}
	err = altOwner.New(acct, curveSpec)
	if err != nil {
		t.Fatal(err)
	}

	utxo := &TXOut{}
	_, err = utxo.Account()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.Account()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}
	ds.DSLinker = &DSLinker{}
	ds.DSLinker.DSPreImage = &DSPreImage{}
	ds.DSLinker.DSPreImage.Owner = &DataStoreOwner{}
	err = ds.DSLinker.DSPreImage.Owner.NewFromOwner(o)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.Account()
	if err != nil {
		t.Fatal(err)
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.Account()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}

	vs.VSPreImage = &VSPreImage{}
	vs.VSPreImage.Owner = &ValueStoreOwner{}
	err = vs.VSPreImage.Owner.NewFromOwner(o)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.Account()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTXOGenericOwner(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	acct := make([]byte, constants.OwnerLen)
	curveSpec := constants.CurveSecp256k1
	o := &Owner{}
	err = o.New(acct, curveSpec)
	if err != nil {
		t.Fatal(err)
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	ds.DSLinker = &DSLinker{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	ds.DSLinker.DSPreImage = &DSPreImage{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}

	ds.DSLinker.DSPreImage.Owner = &DataStoreOwner{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (5)")
	}

	err = ds.DSLinker.DSPreImage.Owner.NewFromOwner(o)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err != nil {
		t.Fatal(err)
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (9)")
	}

	vs.VSPreImage = &VSPreImage{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (10)")
	}

	vs.VSPreImage.Owner = &ValueStoreOwner{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err == nil {
		t.Fatal("Should have raised error (11)")
	}

	err = vs.VSPreImage.Owner.NewFromOwner(o)
	if err != nil {
		t.Fatal(err)
	}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.GenericOwner()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTXOIsDeposit(t *testing.T) {
	utxo := &TXOut{}
	val := utxo.IsDeposit()
	if val {
		t.Fatal("Should be false (1)")
	}

	ds := &DataStore{}
	err := utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	val = utxo.IsDeposit()
	if val {
		t.Fatal("Should be false (3)")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	val = utxo.IsDeposit()
	if val {
		t.Fatal("Should be false (4)")
	}
}

func TestUTXOCannotBeMinedBeforeHeight(t *testing.T) {
	utxo := &TXOut{}
	_, err := utxo.CannotBeMinedBeforeHeight()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ds := &DataStore{}
	err = utxo.NewDataStore(ds)
	if err != nil {
		t.Fatal(err)
	}
	_, err = utxo.CannotBeMinedBeforeHeight()
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
	ds.DSLinker = &DSLinker{}
	ds.DSLinker.DSPreImage = &DSPreImage{}
	ds.DSLinker.DSPreImage.IssuedAt = uint32(0)
	_, err = utxo.CannotBeMinedBeforeHeight()
	if err == nil {
		t.Fatal("Should have raised error (5)")
	}
	ds.DSLinker.DSPreImage.IssuedAt = uint32(3)
	heightDSTrue := 2*constants.EpochLength + 1
	height, err := utxo.CannotBeMinedBeforeHeight()
	if err != nil {
		t.Fatal(err)
	}
	if height != heightDSTrue {
		t.Fatal("Incorrect height for DataStore in CannotBeMinedBeforeHeight")
	}

	vs := &ValueStore{}
	err = utxo.NewValueStore(vs)
	if err != nil {
		t.Fatal(err)
	}
	heightVSTrue := uint32(1)
	height, err = utxo.CannotBeMinedBeforeHeight()
	if err != nil {
		t.Fatal(err)
	}
	if height != heightVSTrue {
		t.Fatal("Incorrect height for ValueStore in CannotBeMinedBeforeHeight")
	}
}
