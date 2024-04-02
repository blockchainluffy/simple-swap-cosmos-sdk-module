#!/usr/bin/env bash


mockgen_cmd="mockgen"
$mockgen_cmd -source=x/simpleswap/expected_keepers/expected_keepers.go -package expected_keepers -destination x/simpleswap/expected_keepers/expected_keepers_mocks.go
