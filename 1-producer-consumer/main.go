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
	"time"
)

func producer(stream Stream, tweets *[]*Tweet, ch chan struct{}) {
	for {
		select {
		default:
			go func() {
				defer func() {
					err := recover()
					if err == nil {
						return
					}

					value, ok := err.(error)
					if !ok {
						return
					}

					if value.Error() == "runtime error: index out of range" {
						return // skip this exception
					}

					panic(recover()) // problems happened
				}()

				tweet, err := stream.Next()
				if err == ErrEOF {
					ch <- *new(struct{}) // tell the channel, this loop should be stopped
				} else {
					*tweets = append(*tweets, tweet)
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
	chEOF := make(chan struct{})
	producer(stream, tweets, chEOF)

	// Consumer
	ch := make(chan string, len(*tweets))
	consumer(*tweets, ch)

	for {
		if cap(ch) == len(*tweets) {
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
