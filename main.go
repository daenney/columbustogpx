package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var epoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

const comma = rune(44)

func main() {
	timezone := flag.String("timezone-offset", "+00:00", "timezone offset as configured on the data logger")
	resetTime := flag.Bool("reset-timestamp", true, "set timestamp to the Unix Epoch")
	pretty := flag.Bool("pretty-print", true, "print XML output with human friendly indentation")

	flag.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(o, "ColumbusToGPX transforms Columbus CSV traces into GPX traces suitable for usage with most mapping programs.\n\n")
		fmt.Fprintf(o, "This command takes a single argument, the path to a file to process. To process multilpe files, call this command multiple times. The output is written to stdout.\n\n")
		fmt.Fprintf(o, "The following options are available:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	input, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gpx := NewGPX()

	sc := bufio.NewScanner(input)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "I") || strings.HasPrefix(line, "#") {
			continue
		}
		elems := strings.FieldsFunc(line, func(c rune) bool { return c == comma })
		if len(elems) < 7 {
			// We expect at least 7 items, Index, Type, Date, Time, Lat, Lon, Ele
			fmt.Fprintln(os.Stderr, "skipping line with less than 7 elements")
			continue
		}
		p := Point{
			Latitude:  latLon(elems[4]),
			Longitude: latLon(elems[5]),
			Elevation: elems[6],
			Timestamp: epoch,
		}
		if !*resetTime {
			date := elems[2]
			ts := elems[3]
			t, err := time.Parse("060102150405-07:00", date+ts+*timezone)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to parse date: %s and time: %s with offset: %s, using Unix epoch instead: %s\n", date, ts, *timezone, err)
			} else {
				p.Timestamp = t.UTC()
			}
		}

		switch cat := elems[1]; cat {
		case "T":
			gpx.AddTrackSegment(p)
		case "C":
			gpx.AddPOI(p)
		default:
			fmt.Fprintf(os.Stderr, "point %s excluded due to unknown type: %s\n", elems[0], cat)
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		input.Close()
		os.Exit(1)
	}

	e := xml.NewEncoder(os.Stdout)
	if *pretty {
		e.Indent("", "\t")
	}
	fmt.Fprintln(os.Stdout, `<?xml version="1.0" encoding="UTF-8"?>`)
	if err := e.Encode(gpx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		input.Close()
		os.Exit(1)
	}
	input.Close()
}

func latLon(s string) string {
	idx := len(s) - 1
	last := s[idx:]

	switch last {
	case "N", "E":
		return s[:idx]
	case "W", "S":
		return "-" + s[:idx]
	default:
		panic("latitude/longitude without cardinal direction")
	}
}
