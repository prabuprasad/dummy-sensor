package main

 import (
    "encoding/json"
    "math"
 	"context"
 	"log"
 	"time"
 	"fmt"
 	"os"
 	"strings"
    "github.com/pkg/errors"
 	"pack.ag/amqp"
 )

func echo(args []string) error {
   if len(args) < 2 {
       return errors.New("no message to echo")
   }
   _, err := fmt.Println(strings.Join(args[1:], " "))
   return err
}
type Signal struct {
	Timestamp time.Time `json:"timestamp"`
	Value float64 `json:"value"`
	Sensor string  `json:"sensor"`
}

 func main() {
 	// Create client
 	amqpEndpoint := os.Getenv("endpoint")
    if amqpEndpoint == "" {
   		amqpEndpoint = "amqp://localhost:5672"
    }
 	client, err := amqp.Dial(amqpEndpoint,
  	//	amqp.ConnSASLPlain("cornerstone", "LikeABosch123"),
 	)
    channel := os.Getenv("channel")
   	if channel == "" {
   		channel = "radium3"
   	}
   	sensor := os.Getenv("sensor")
   	if sensor == "" {
       		sensor = "sensor99"
    }
 	if err != nil {
 	    fmt.Fprintf(os.Stderr, "%+v\n", err)
 		log.Fatal("Dialing AMQP server:", err)
 	}
 	defer client.Close()

 	// Open a session
 	session, err := client.NewSession()
 	if err != nil {
 		log.Fatal("Creating AMQP session:", err)
 	}

 	ctx := context.Background()

	// Send message
 	function := os.Getenv("function")

    ticker := time.NewTicker(time.Millisecond*20)
    tickerOutput := time.NewTicker(time.Second*10)
    counter := 0

    for {
    // Create a sender
 	sender, err := session.NewSender(
 		amqp.LinkTargetAddress(channel),
 	)
 	if err != nil {
 		log.Fatal("Creating sender link:", err)
 	}

    ctx, cancel := context.WithTimeout(ctx, 1000*time.Second)

    	select {
    	    case curTime := <- ticker.C:
            var value float64
        switch function {
        // linear
        case "linear":
        		value = (float64(curTime.UTC().UnixNano() % (int64(time.Second) * 10)) / float64(time.Second))/5 -1
        case "bool":
        	if  curTime.UTC().UnixNano() % (int64(time.Second) * 10) > int64(time.Second*5) {
        		value = 1
        		} else {
        					value = -1
        				}
        				// sine function for the default
        			default:
        				value = math.Sin(float64(curTime.UTC().UnixNano())/ float64(time.Second) *10 / ( 2.0 * math.Pi))
        			}


        			signal := &Signal{
        			    Sensor: sensor,
        				Timestamp: curTime,
        				Value:     value,
        			}
        			raw, err := json.Marshal(signal)
        			if err != nil {
        				log.Printf("Error marshaling to json: %+v", err)
        				continue
        			}
        			msg := amqp.NewMessage(raw)
                    err = sender.Send(ctx, msg)

//         			err = sender.Send(ctx, amqp.NewMessage([]byte("Hello")))

        			if err != nil {
                      log.Printf("Creating sender link:", err)
                    }
        			counter++

        		case <- tickerOutput.C:
        			log.Printf("Sent %v sensor values in 10 seconds", counter)
        			counter = 0
        		}
 		//end sender
 		if err != nil {
 			log.Fatal("Sending message:", err)
 		}

  		sender.Close(ctx)
 		cancel()
 	    }

}
