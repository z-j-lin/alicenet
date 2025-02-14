#! /bin/bash
# DO NOT EDIT THIS FILE! THIS FILE CONTAINS GENERATED CODE
# CHANGES MADE TO THIS FILE WILL BE LOST
#
# npx hardhat go-go-gen --pkg bindings --in bindings/bindings-artifacts --out bindings 
abigen --abi bindings/bindings-artifacts/AToken.json --pkg bindings --type AToken --out bindings/AToken.go
abigen --abi bindings/bindings-artifacts/ATokenBurner.json --pkg bindings --type ATokenBurner --out bindings/ATokenBurner.go
abigen --abi bindings/bindings-artifacts/ATokenMinter.json --pkg bindings --type ATokenMinter --out bindings/ATokenMinter.go
abigen --abi bindings/bindings-artifacts/AliceNetFactory.json --pkg bindings --type AliceNetFactory --out bindings/AliceNetFactory.go
abigen --abi bindings/bindings-artifacts/BToken.json --pkg bindings --type BToken --out bindings/BToken.go
abigen --abi bindings/bindings-artifacts/BTokenErrors.json --pkg bindings --type BTokenErrors --out bindings/BTokenErrors.go
abigen --abi bindings/bindings-artifacts/ETHDKG.json --pkg bindings --type ETHDKG --out bindings/ETHDKG.go
abigen --abi bindings/bindings-artifacts/ETHDKGErrors.json --pkg bindings --type ETHDKGErrors --out bindings/ETHDKGErrors.go
abigen --abi bindings/bindings-artifacts/Governance.json --pkg bindings --type Governance --out bindings/Governance.go
abigen --abi bindings/bindings-artifacts/GovernanceErrors.json --pkg bindings --type GovernanceErrors --out bindings/GovernanceErrors.go
abigen --abi bindings/bindings-artifacts/PublicStaking.json --pkg bindings --type PublicStaking --out bindings/PublicStaking.go
abigen --abi bindings/bindings-artifacts/Snapshots.json --pkg bindings --type Snapshots --out bindings/Snapshots.go
abigen --abi bindings/bindings-artifacts/SnapshotsErrors.json --pkg bindings --type SnapshotsErrors --out bindings/SnapshotsErrors.go
abigen --abi bindings/bindings-artifacts/ValidatorPool.json --pkg bindings --type ValidatorPool --out bindings/ValidatorPool.go
abigen --abi bindings/bindings-artifacts/ValidatorPoolErrors.json --pkg bindings --type ValidatorPoolErrors --out bindings/ValidatorPoolErrors.go
abigen --abi bindings/bindings-artifacts/ValidatorStaking.json --pkg bindings --type ValidatorStaking --out bindings/ValidatorStaking.go
