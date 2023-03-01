vegeta attack -targets vegeta.txt -header Authorization:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxIiwiZXhwIjoxNjc3NzkxNTAyfQ.CDTEtPX1gaIH1-F6hdTuAw6-HSWF_6LtBPGzipEgItg" -name=1500qps -rate=1500 -duration=600s > results.1500qps.bin;cat results.1500qps.bin | vegeta plot > plot.1500qps.html;cat results.1500qps.bin | vegeta report -type=text

