package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/require"

	"github.com/NibiruChain/nibiru/x/common/testutil/testapp"
	"github.com/NibiruChain/nibiru/x/epochs/keeper"

	epochstypes "github.com/NibiruChain/nibiru/x/epochs/types"
)

func TestQueryEpochInfos(t *testing.T) {
	nibiruApp, ctx := testapp.NewNibiruTestAppAndContext()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, nibiruApp.InterfaceRegistry())
	epochstypes.RegisterQueryServer(queryHelper, keeper.NewQuerier(nibiruApp.EpochsKeeper))
	queryClient := epochstypes.NewQueryClient(queryHelper)

	epochInfos := nibiruApp.EpochsKeeper.AllEpochInfos(ctx)
	chainStartTime := epochInfos[0].StartTime
	errMsg := fmt.Sprintf("epochInfos: %v\n", epochInfos)

	// Invalid param
	epochInfosResponse, err := queryClient.EpochInfos(
		gocontext.Background(), &epochstypes.QueryEpochsInfoRequest{},
	)
	require.NoError(t, err, errMsg)
	require.Len(t, epochInfosResponse.Epochs, 3)

	// check if EpochInfos are correct
	require.Equal(t, epochInfosResponse.Epochs[0].StartTime, chainStartTime, errMsg)
	require.Equal(t, epochInfosResponse.Epochs[0].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, epochInfosResponse.Epochs[0].Identifier, "15 min")
	require.Equal(t, epochInfosResponse.Epochs[0].Duration, time.Minute*15)
	require.Equal(t, epochInfosResponse.Epochs[0].CurrentEpoch, uint64(0))
	require.Equal(t, epochInfosResponse.Epochs[0].EpochCountingStarted, false)

	require.Equal(t, epochInfosResponse.Epochs[1].StartTime, chainStartTime)
	require.Equal(t, epochInfosResponse.Epochs[1].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, epochInfosResponse.Epochs[1].Identifier, "30 min")
	require.Equal(t, epochInfosResponse.Epochs[1].Duration, time.Minute*30)
	require.Equal(t, epochInfosResponse.Epochs[1].CurrentEpoch, uint64(0))
	require.Equal(t, epochInfosResponse.Epochs[1].EpochCountingStarted, false)
}

func TestCurrentEpochQuery(t *testing.T) {
	nibiruApp, ctx := testapp.NewNibiruTestAppAndContext()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, nibiruApp.InterfaceRegistry())
	epochstypes.RegisterQueryServer(queryHelper, keeper.NewQuerier(nibiruApp.EpochsKeeper))
	queryClient := epochstypes.NewQueryClient(queryHelper)

	// Valid epoch
	epochInfosResponse, err := queryClient.CurrentEpoch(gocontext.Background(), &epochstypes.QueryCurrentEpochRequest{Identifier: "15 min"})
	require.NoError(t, err)
	require.Equal(t, epochInfosResponse.CurrentEpoch, uint64(0))

	// Invalid epoch
	_, err = queryClient.CurrentEpoch(gocontext.Background(), &epochstypes.QueryCurrentEpochRequest{Identifier: "invalid epoch"})
	require.Error(t, err)
}
