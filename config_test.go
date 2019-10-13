package main

import "testing"

func Benchmark_Config(b *testing.B) {
	ConfigDirectory = "./configs"
	initConfigs()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		switch n % 5 {
		case 0:
			v := configViper["core"].GetInt("port")
			if v != 8443 {
				b.Fail()
			}
		case 1:
			v := configViper["picking"].GetString("qa-int.us.r1.capability")
			if v != "https://qa-us-r1-picking.com" {
				b.Fail()
			}
		case 2:
			v := configViper["picking"].GetString("qa-int.uk.capability")
			if v != "https://qa-uk-picking.com" {
				b.Fail()
			}
		case 3:
			v := configViper["packing"].GetBool("serialNumberEnabled")
			if !v {
				b.Fail()
			}
		case 4:
			v := configViper["packing"].GetBool("qa-int.us.r2.1.5594.serialNumberEnabled")
			if v {
				b.Fail()
			}
		}
	}
}
