package kafkabuilder

import (
	"github.com/segmentio/kafka-go"
)

func NewReader(cfg *kafka.ReaderConfig) *kafka.Reader {
	reader := kafka.NewReader(*cfg)
	return reader
}

func NewWriter(cfg *kafka.WriterConfig) *kafka.Writer {
	writer := kafka.NewWriter(*cfg)
	return writer
}

// func Auth(algo *scram.Algorithm, user string, pass string) sasl.Mechanism {
// 	if algo == nil {
// 		// plain
// 		return plain.Mechanism{
// 			Username: "user",
// 			Password: "pass",
// 		}
// 	}
// 	// scram
// 	mechanism, _ := scram.Mechanism(*algo, user, pass)
// 	return mechanism
// }
