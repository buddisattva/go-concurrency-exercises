//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

type tweetsWithMutex struct {
	tweets []*Tweet
	mux    sync.Mutex
}

func producer(stream Stream, tweetsWithMutex *tweetsWithMutex, ch chan struct{}) {
	for {
		select {
		default:
			go func() {
				defer tweetsWithMutex.mux.Unlock()
				tweetsWithMutex.mux.Lock()

				tweet, err := stream.Next()
				if err == ErrEOF {
					ch <- *new(struct{}) // tell the channel, this loop should be stopped
				} else {
					tweets := tweetsWithMutex.tweets
					tweets = append(tweets, tweet)
					tweetsWithMutex.tweets = tweets
				}
			}()
		case <-ch:
			return
		}
	}
}

func consumer(tweets []*Tweet, ch chan string) {
	for _, t := range tweets {
		cloned := new(Tweet)
		*cloned = *t
		go consumeTweet(*cloned, ch)
	}
}

func consumeTweet(t Tweet, ch chan string) {
	if t.IsTalkingAboutGo() {
		ch <- t.Username + "\ttweets about golang"
	} else {
		ch <- t.Username + "\tdoes not tweet about golang"
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweets := new([]*Tweet)
	tweetsWithMutex := tweetsWithMutex{
		tweets: *tweets,
	}
	chEOF := make(chan struct{})
	producer(stream, &tweetsWithMutex, chEOF)

	// Consumer
	ch := make(chan string, len(tweetsWithMutex.tweets))
	consumer(tweetsWithMutex.tweets, ch)

	for {
		if cap(ch) == len(tweetsWithMutex.tweets) {
			break
		}
		time.Sleep(time.Nanosecond)
	}

	for msg := range ch {
		fmt.Println(msg)
		if len(ch) < 1 {
			break
		}
	}

	fmt.Printf("Process took %s\n", time.Since(start))
}
