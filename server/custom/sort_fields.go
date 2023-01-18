package custom

import (
	"fmt"
	ps "partitions_searcher/proto/partitions_searcher"
	"partitions_searcher/settings"
)

func GetDifference(r1 *ps.Record, r2 *ps.Record, direction bool) bool {
	if direction {
		return r1.Id <= r2.Id
	}
	return r1.Id >= r2.Id
}

func GetEventByPartition(partition string, query string) (records []ps.Record, err error) {
	err = settings.Db.Raw(fmt.Sprintf("select id, title, description from %s as e %s", partition, query)).Scan(&records).Error
	return
}
