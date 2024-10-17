package numstat

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"sort"
	"strconv"
)

type Stats struct {
	Entries     int
	UniqEntries int
	Buckets     []Bucket
	Percentiles []Percentile
	Min         float64
	Max         float64
	Sum         *big.Float
	Avg         float64
	Var         float64
	StdDev      float64
	ER680       Range
	ER950       Range
	ER997       Range
}

type Range struct {
	Start float64
	End   float64
}

type Bucket struct {
	Range
	Entries int
}

type Percentile struct {
	Percentile float64
	Value      float64
}

func CalcStats(data []float64, numBuckets int) Stats {
	if len(data) == 0 {
		return Stats{}
	}

	stats := Stats{
		Entries: len(data),
		Min:     data[0],
		Max:     data[len(data)-1],
	}

	//
	// First iteration for Sum, Average and Uniq values
	//
	lengthF := float64(len(data))
	stats.Sum = big.NewFloat(0)
	var lastDatum float64
	for idx, datum := range data {
		// sum
		stats.Sum.Add(stats.Sum, big.NewFloat(datum))

		// uniqueness
		if lastDatum != datum || idx == 0 {
			stats.UniqEntries++
		}
		lastDatum = datum
	}
	// average
	avg := big.NewFloat(0).Copy(stats.Sum)
	avg.Quo(stats.Sum, big.NewFloat(lengthF))
	stats.Avg, _ = avg.Float64()

	//
	// Second iteration for Bucketing, Variance and Standard Deviation
	//
	if numBuckets > stats.UniqEntries {
		numBuckets = stats.UniqEntries
	}
	stats.Buckets = make([]Bucket, numBuckets)
	numBucketsF := float64(numBuckets)
	bucketRange := (stats.Max - stats.Min) / numBucketsF
	for idx := range stats.Buckets {
		stats.Buckets[idx].Start = stats.Min + float64(idx)*bucketRange
		stats.Buckets[idx].End = stats.Min + float64(idx+1)*bucketRange
	}
	for _, datum := range data {
		// bucketing
		bucketIdx := int((datum - stats.Min) / bucketRange)
		if bucketIdx == numBuckets {
			bucketIdx = numBuckets - 1
		}
		stats.Buckets[bucketIdx].Entries++

		// variance
		stats.Var += math.Pow(datum-stats.Avg, 2) / lengthF
	}
	// standard deviation
	stats.StdDev = math.Sqrt(stats.Var)

	//
	// Empirical Rule https://en.wikipedia.org/wiki/68%E2%80%9395%E2%80%9399.7_rule
	//
	stats.ER680 = Range{Start: stats.Avg - stats.StdDev*1.0, End: stats.Avg + stats.StdDev*1.0}
	stats.ER950 = Range{Start: stats.Avg - stats.StdDev*2.0, End: stats.Avg + stats.StdDev*2.0}
	stats.ER997 = Range{Start: stats.Avg - stats.StdDev*3.0, End: stats.Avg + stats.StdDev*3.0}

	//
	// Percentiles
	//
	stats.Percentiles = []Percentile{}
	addPct := func(pct float64, value float64) {
		entry := Percentile{
			Percentile: pct,
			Value:      value,
		}
		stats.Percentiles = append(stats.Percentiles, entry)
	}
	// linear progression (+1) from pct1 to pct10 and pct90 to pct99
	for i := 1; i < 10; i++ {
		lowIdx := float64(i) / 100 * (lengthF - 1)
		highIdx := lengthF - 1 - lowIdx
		addPct(float64(i), data[int(math.Round(lowIdx))])
		addPct(100-float64(i), data[int(math.Round(highIdx))])
	}
	// linear progression (+10) from pct10 to pct90
	for i := 10; i <= 90; i += 10 {
		idx := float64(i) / 100 * (lengthF - 1)
		addPct(float64(i), data[int(math.Round(idx))])
	}
	// logaritmic progression from pct0.0001 to pct0.1
	// and from pct99.9 to pct99.9999
	for digits := 1; digits <= 4; digits++ {
		pct1 := 1.0 / math.Pow(10, float64(digits))
		pct99 := 100.0 - pct1
		addPct(pct1, data[int(pct1/100.0*(lengthF-1))])
		addPct(pct99, data[int(pct99/100.0*(lengthF-1))])
	}
	sort.Slice(stats.Percentiles, func(i, j int) bool {
		return stats.Percentiles[i].Percentile < stats.Percentiles[j].Percentile
	})

	return stats
}

func (s *Stats) Print(writer io.Writer) error {
	var err error
	fmt.Fprintf(writer, "entries\t%d\n", s.Entries)
	fmt.Fprintf(writer, "uniq\t%d\n", s.UniqEntries)
	fmt.Fprintf(writer, "min\t%.5f\n", s.Min)
	fmt.Fprintf(writer, "max\t%.5f\n", s.Max)
	fmt.Fprintf(writer, "avg\t%.5f\n", s.Avg)
	fmt.Fprintf(writer, "sum\t%s\n", s.Sum.String())
	fmt.Fprintf(writer, "varce\t%.5f\n", s.Var)
	fmt.Fprintf(writer, "stddev\t%.5f\n", s.StdDev)
	fmt.Fprintf(writer, "er68.0\t%.5f ~ %.5f\n", s.ER680.Start, s.ER680.End)
	fmt.Fprintf(writer, "er95.0\t%.5f ~ %.5f\n", s.ER950.Start, s.ER950.End)
	fmt.Fprintf(writer, "er99.7\t%.5f ~ %.5f\n", s.ER997.Start, s.ER997.End)

	for _, pct := range s.Percentiles {
		fmt.Fprintf(writer, "pct%-5g\t%.5f\n", pct.Percentile, pct.Value)
	}

	digits := int(math.Log10(float64(s.Entries))) + 1
	fmtStr := "bkt%02d\t%.5f ~ %.5f\tentries %-" + strconv.Itoa(digits) + "d\t%0.2f%%\t%s\n"
	for idx, bucket := range s.Buckets {
		percentage := float64(bucket.Entries*100) / float64(s.Entries)
		_, err = fmt.Fprintf(
			writer,
			fmtStr,
			idx+1, bucket.Start, bucket.End, bucket.Entries, percentage,
			Bar(percentage, 100.0, true),
		)
	}
	return err
}
