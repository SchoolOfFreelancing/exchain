package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/okex/exchain/x/slashing/internal/types"
)

// Simulation parameter constants
const (
	SignedBlocksWindow      = "signed_blocks_window"
	MinSignedPerWindow      = "min_signed_per_window"
	DowntimeJailDuration    = "downtime_jail_duration"
	SlashFractionDoubleSign = "slash_fraction_double_sign"
	SlashFractionDowntime   = "slash_fraction_downtime"
)

// GenSignedBlocksWindow randomized SignedBlocksWindow
func GenSignedBlocksWindow(r *rand.Rand) int64 {
	return int64(simulation.RandIntBetween(r, 10, 1000))
}

// GenMinSignedPerWindow randomized MinSignedPerWindow
func GenMinSignedPerWindow(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(10)), 1)
}

// GenDowntimeJailDuration randomized DowntimeJailDuration
func GenDowntimeJailDuration(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24)) * time.Second
}

// GenSlashFractionDoubleSign randomized SlashFractionDoubleSign
func GenSlashFractionDoubleSign(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(1).Quo(sdk.NewDec(int64(r.Intn(50) + 1)))
}

// GenSlashFractionDowntime randomized SlashFractionDowntime
func GenSlashFractionDowntime(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(1).Quo(sdk.NewDec(int64(r.Intn(200) + 1)))
}

// RandomizedGenState generates a random GenesisState for slashing
func RandomizedGenState(simState *module.SimulationState) {
	var signedBlocksWindow int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SignedBlocksWindow, &signedBlocksWindow, simState.Rand,
		func(r *rand.Rand) { signedBlocksWindow = GenSignedBlocksWindow(r) },
	)

	var minSignedPerWindow sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinSignedPerWindow, &minSignedPerWindow, simState.Rand,
		func(r *rand.Rand) { minSignedPerWindow = GenMinSignedPerWindow(r) },
	)

	var downtimeJailDuration time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, DowntimeJailDuration, &downtimeJailDuration, simState.Rand,
		func(r *rand.Rand) { downtimeJailDuration = GenDowntimeJailDuration(r) },
	)

	var slashFractionDoubleSign sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SlashFractionDoubleSign, &slashFractionDoubleSign, simState.Rand,
		func(r *rand.Rand) { slashFractionDoubleSign = GenSlashFractionDoubleSign(r) },
	)

	var slashFractionDowntime sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SlashFractionDowntime, &slashFractionDowntime, simState.Rand,
		func(r *rand.Rand) { slashFractionDowntime = GenSlashFractionDowntime(r) },
	)

	params := types.NewParams(
		signedBlocksWindow, minSignedPerWindow, downtimeJailDuration,
		slashFractionDoubleSign, slashFractionDowntime,
	)

	slashingGenesis := types.NewGenesisState(params, nil, nil)

	fmt.Printf("Selected randomly generated slashing parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, slashingGenesis.Params))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(slashingGenesis)
}
