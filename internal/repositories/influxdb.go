package repositories

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	"time"
)

type InfluxDB struct {
	w api.WriteAPI
}

func (s InfluxDB) Send(r *proto.Metrics) error {
	for _, m := range r.Data() {
		tags := map[string]string{}
		for i := 0; i < len(m.Tags)-1; i += 2 {
			tags[m.Tags[i]] = m.Tags[i+1]
		}
		for i := 0; i < len(m.DefaultTags)-1; i += 2 {
			tags[m.DefaultTags[i]] = m.DefaultTags[i+1]
		}
		s.w.WritePoint(influxdb2.NewPoint(m.Name, tags, map[string]interface{}{
			`min`:     m.Values.Min,
			`max`:     m.Values.Max,
			`count`:   m.Values.Count,
			`sum`:     m.Values.Sum,
			`first`:   m.Values.First,
			`last`:    m.Values.Last,
			`average`: m.Values.Sum / float64(m.Values.Count),
		}, time.Unix(0, m.Time)))
	}
	return nil
}

func NewInfluxDB(v *viper.Viper) Repo {

	client := influxdb2.NewClientWithOptions(v.GetString(`INFLUXDB_URL`), v.GetString(`INFLUXDB_TOKEN`), influxdb2.DefaultOptions().SetFlushInterval(1000))
	writeAPI := client.WriteAPI(v.GetString(`INFLUXDB_ORG`), v.GetString(`INFLUXDB_BUCKET`))
	s := &InfluxDB{
		w: writeAPI,
	}
	go func() {
		for err := range writeAPI.Errors() {
			fmt.Println(err)
		}
	}()

	return s
}
