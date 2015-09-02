package kafka

import ( 
        "github.com/Shopify/sarama"
        "encoding/json" 
        "reflect"
        "log"
)


type Producer struct {
  syncProducer sarama.SyncProducer
}

// Constructor
func NewProducer(brokerList []string) *Producer {
  config := sarama.NewConfig()
  config.Producer.RequiredAcks = sarama.WaitForAll 
  config.Producer.Retry.Max = 10    
  syncProducer, err := sarama.NewSyncProducer(brokerList, config)
  if err != nil {
    log.Fatalln("Failed to start Sarama producer:", err)
    panic(err)
  }

  return &Producer{ syncProducer : syncProducer }
}

func (this *Producer) SendEventToTopic( event interface{}, topic string ) error {
  
  // marshal event
  json, err := json.Marshal(event)
  
  if err != nil {
    return err
  }

  log.Println(string(json))

  // send event
  _, _, err = this.syncProducer.SendMessage(&sarama.ProducerMessage {
      Topic: topic,
      Value: sarama.StringEncoder(reflect.TypeOf(event).Name() + "," + string(json)),
    })

  if err != nil {
    return err
  }

  log.Println("event sent")
  return nil
}



// defer func() {
//   if err := Producer.Close(); err != nil {
//     log.Println("Failed to shut down data collector cleanly", err)
//   }
// }()