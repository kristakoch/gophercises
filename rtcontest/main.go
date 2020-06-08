package main

import (
	"context"
	"log"
)

const (
	tweet = "1173117502848741376" // chrissyteigen tweet
	user  = "eedrk"               // good twitter guy
)

func main() {
	ctx := context.Background()

	t, err := New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Run timed rt contest.
	// Pick 3 winners. Get data every 5 minutes for 1 hour.
	t.TimedRTContest(ctx, tweet, 3, 5, 1)

	// Run rt contest.
	// if err = t.GetAndStoreRTs(ctx, tweet); err != nil {
	// 	log.Fatal(err)
	// }

	// if err = t.RTContestWinners(ctx, tweet, 1); err != nil {
	// 	log.Fatal(err)
	// }

	// // Log tweets by user.
	// if err = t.LogUserTweets(ctx, user, 10); err != nil {
	// 	log.Fatal(err)
	// }

	// // Log favs by user.
	// if err = t.LogUserFavs(ctx, user, 10); err != nil {
	// 	log.Fatal(err)
	// }
}
