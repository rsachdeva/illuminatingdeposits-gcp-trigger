//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_analytics

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// q := client.Query(
// 	`select
// 		  name
// 		from
// 		  bigquery-public-data.usa_names.usa_1910_2013
// 		where
// 		  state = "TX"
// 		order by
// 		    name
// 		limit
// 		  3`)
// // "select name from `bigquery-public-data.usa_names.usa_1910_2013` " +
// // "WHERE state = \"TX\" " +
// // "LIMIT 3"
// // Location must match that of the dataset(depositsByDate) referenced in the query.
// q.Location = "US"
// // Run the query and process the returned row iterator.
// iter, err := q.Read(ctx)
// if err != nil {
// 	return fmt.Errorf("query.Read(): %w", err)
// }
// for {
// 	var row []bigquery.Value
// 	err := iter.Next(&row)
// 	if err == iterator.Done {
// 		break
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("row is", row)
// }
// return nil

// q := client.Query(
// 	`select word, word_count
//     from ` + "`bigquery-public-data.samples.shakespeare`" + `
//     WHERE corpus = @corpus
//     AND word_count >= @min_word_count
//     ORDER BY word_count DESC;`)
// q.Parameters = []bigquery.QueryParameter{
// 	{
// 		Name:  "corpus",
// 		Value: "romeoandjuliet",
// 	},
// 	{
// 		Name:  "min_word_count",
// 		Value: 250,
// 	},
// }

// const depositsByDate = `select
//
//	  dc.delta,
//	  format_date("%a,%x", dc.created_at) depositsByDate
//	from
//	  illuminatingdeposits-gcp.gcfdeltaanalytics.delta_calculations dc
//	order by
//	  dc.created_at desc
//	limit
//	  3;`
func queryDepositHighestDelta(ctx context.Context, writer io.Writer) error {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %writer", err)
	}
	defer client.Close()

	const depositsByDate = `select
			cast(dc.delta as string) delta,
			format_date("%a,%x", dc.created_at) depositsByDate
		from
			illuminatingdeposits-gcp.gcfdeltaanalytics.delta_calculations dc
		order by
			dc.created_at desc
		limit
			10;`

	q := client.Query(depositsByDate)

	// Run the query and process the returned row iterator.
	iter, err := q.Read(ctx)
	if err != nil {
		return fmt.Errorf("query.Read(): %writer", err)
	}

	fmt.Fprintln(writer, []bigquery.Value{"delta", "depositsByDate"})

	for {
		var row []bigquery.Value

		err := iter.Next(&row)
		// if err == iterator.Done {
		// 	break
		// }
		log.Println("bigquery errors.Is(err, iterator.Done) is", errors.Is(err, iterator.Done))

		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return fmt.Errorf("queryDepositHighestDelta iterator.Next(): %writer", err)
		}
		// log.Printf("row is %#v", row)
		// log.Printf("len(row) is %v", len(row))
		// log.Println("row[0] is", row[0])
		// log.Printf("row[0] is %#v\n", row[0])
		// ratioValue := row[0].(*big.Rat)
		// log.Println("row[0] in number calculated", ratioValue.FloatString(2))
		// log.Println("row[1] is", row[1])
		// log.Printf("row[1] is %#v\n", row[1])
		fmt.Fprintln(writer, row)
	}

	return nil
}
