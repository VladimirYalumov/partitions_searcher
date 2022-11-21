package main

import (
	"fmt"
	"partitions_searcher/example/models"
	"partitions_searcher/settings"
	"strconv"
	"time"
)

const StartYear = 2022 // year from which partitions should be created
const MonthCount = 24  // count of partitions

type CheckTableStruct struct {
	TableName string `json:"table_name"`
}

type CheckTriggerStruct struct {
	Count int `json:"count"`
}

func main() {
	settings.Init("../../config.yml")
	settings.Db.AutoMigrate(&models.Task{})
	generateTasksPartitions()

	// create 2 records in every partition
	for i := 1; i < 12; i++ {
		task := models.Task{
			Title:       "test " + strconv.Itoa(i),
			Description: "test description",
			CreatedAt:   time.Date(2022, time.Month(i), 1, 0, 0, 0, 0, time.UTC),
		}
		_ = task.Create()
	}
	for i := 1; i < 12; i++ {
		task := models.Task{
			Title:       "test " + strconv.Itoa(i),
			Description: "test description",
			CreatedAt:   time.Date(2022, time.Month(i), 1, 0, 0, 0, 0, time.UTC),
		}
		_ = task.Create()
	}

	// create random records in every partition
	for i := 1; i < 12; i++ {
		task := models.Task{
			Title:       "qwerty",
			Description: "random description",
			CreatedAt:   time.Date(2022, time.Month(i), 1, 0, 0, 0, 0, time.UTC),
		}
		_ = task.Create()
	}
}

func DateStart(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func DateEnd(year, month, day int) time.Time {
	date := time.Date(year, time.Month(month+1), day, 0, 0, 0, 0, time.UTC)
	return date.Add(-1 * time.Millisecond)
}

func GetTaskPartitionName(year, month int) string {
	if month < 10 {
		return fmt.Sprintf("tasks_%d_0%d", year, month)
	}
	return fmt.Sprintf("tasks_%d_%d", year, month)
}

func generateTasksPartitions() {
	var year = StartYear
	var month int
	var result []CheckTableStruct
	var triggerResult CheckTriggerStruct
	var err error

	for i := 0; i < MonthCount; i++ {
		modification := i / 12
		month = (i + 1) - (modification)*12
		year = modification + StartYear
		settings.Db.Raw("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME=?;", GetTaskPartitionName(year, month)).Scan(&result)
		if len(result) != 0 {
			continue
		}
		partitionStartDate := DateStart(year, month, 1)
		partitionEndDate := DateEnd(year, month, 1)
		err = settings.Db.Exec(fmt.Sprintf("create table \"%s\" ( like tasks including all );", GetTaskPartitionName(year, month))).Error
		if err != nil {
			panic(err)
		}
		err = settings.Db.Exec(fmt.Sprintf("alter table \"%s\" add constraint partition_check check (created_at between '%v' and '%v');",
			GetTaskPartitionName(year, month),
			partitionStartDate.Format(`2006-01-02 15:04:05.000 -0700`),
			partitionEndDate.Format(`2006-01-02 15:04:05.000 -0700`))).Error
		if err != nil {
			panic(err)
		}
		err = settings.Db.Exec(fmt.Sprintf("alter table \"%s\" inherit tasks", GetTaskPartitionName(year, month))).Error
		if err != nil {
			panic(err)
		}
	}

	settings.Db.Raw("select count(*) as count " +
		"from pg_trigger " +
		"where not tgisinternal " +
		"and tgrelid = 'tasks'::regclass and tgname = 'partition_tasks';").Scan(&triggerResult)

	if triggerResult.Count == 0 {
		err = settings.Db.Exec("create or replace function partition_for_tasks() returns trigger as " +
			"$$ " +
			"DECLARE " +
			"v_parition_name text; " +
			"BEGIN " +
			"v_parition_name := format( 'tasks_%s', to_char(NEW.created_at, 'YYYY_MM') ); " +
			"execute 'INSERT INTO ' || v_parition_name || ' VALUES ( ($1).* )' USING NEW; " +
			"return NEW; " +
			"END; " +
			"$$ " +
			"language plpgsql; ").Error

		if err != nil {
			panic(err)
		}

		err = settings.Db.Exec("create trigger partition_tasks before insert on tasks for each row execute procedure partition_for_tasks();").Error

		if err != nil {
			panic(err)
		}
	}

	settings.Db.Raw("select count(*) as count " +
		"from pg_trigger " +
		"where not tgisinternal " +
		"and tgrelid = 'tasks'::regclass and tgname = 'delete_record_after_insert_tasks';").Scan(&triggerResult)

	if triggerResult.Count == 0 {
		err = settings.Db.Exec("create or replace function delete_record_after_insert_tasks() returns trigger as $$ " +
			"DECLARE delete_sql_query text; " +
			"BEGIN " +
			"delete_sql_query = format( 'DELETE FROM ONLY tasks WHERE id = %s', NEW.id); " +
			"execute delete_sql_query; " +
			"return NEW; " +
			"END; " +
			"$$ " +
			"language plpgsql; ").Error

		if err != nil {
			panic(err)
		}

		err = settings.Db.Exec("create trigger delete_record_after_insert_tasks after insert on tasks for each row execute procedure delete_record_after_insert_tasks();").Error

		if err != nil {
			panic(err)
		}
	}
}
