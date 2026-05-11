package types

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var MintedSatsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "tx_msg_confirm_btc_deposit_total",
	Help: "Total number of minted sats",
})

var OracleBlockGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "oracle_block_number",
	Help: "Block number that each oracle is voting on",
}, []string{"oracle"})
