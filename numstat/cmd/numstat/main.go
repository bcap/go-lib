package main

import (
	"context"
	"flag"
	"os"

	"github.com/bcap/go-lib/numstat"
)

func main() {
	var numBuckets int
	flag.IntVar(&numBuckets, "buckets", 20, "how many buckets to bucket the data into")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data, err := numstat.LoadData(ctx, os.Stdin)
	panicOnErr(err)

	stats := numstat.CalcStats(data, numBuckets)
	panicOnErr(stats.Print(os.Stdout))
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
