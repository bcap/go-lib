package numstat

import (
	"bufio"
	"context"
	"io"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

func LoadData(ctx context.Context, reader io.Reader) ([]float64, error) {
	group, ctx := errgroup.WithContext(ctx)
	c := make(chan string, 10)

	bufReader := bufio.NewReader(reader)
	data := make([]float64, 0, 1024)

	// producer: read lines
	group.Go(func() error {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			line, err := bufReader.ReadString('\n')
			if err == io.EOF && line == "" {
				return nil
			}
			if err != nil && err != io.EOF {
				return err
			}
			select {
			case c <- line:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	// consumer: process data
	group.Go(func() error {
		for {
			var line string
			var read bool
			select {
			case line, read = <-c:
				if !read {
					return nil
				}
			case <-ctx.Done():
				return ctx.Err()
			}

			line = strings.TrimSpace(line)
			v, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return err
			}

			data = append(data, v)
		}
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}
	sort.Float64s(data)
	return data, nil
}
