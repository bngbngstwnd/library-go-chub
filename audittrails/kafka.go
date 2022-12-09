package audittrails

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"time"

	"github.com/bngbngstwnd/library-go-chub/config"
	"github.com/bngbngstwnd/library-go-chub/constant"
	"github.com/bngbngstwnd/library-go-chub/model/entity"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riferrei/srclient"

	"gorm.io/gorm"
)

type AuditTrails interface {
	Produce(entity.AuditTrail) error
	Store(entity.AuditTrail) error
	Middleware() gin.HandlerFunc
}

type auditTrails struct {
	serviceName string
	kafka       KafkaConfig
	producer    *kafka.Producer
	schema      *srclient.Schema
	mysql       *gorm.DB
}

type KafkaConfig struct {
	Brokers        string
	SchemaRegistry string
	SASLUsername   string
	SASLPassword   string
}

func NewAuditTrailsConn(serviceName string, conf KafkaConfig, mysql *gorm.DB) (AuditTrails, error) {

	producer, err := initProducer(conf.Brokers, conf.SASLUsername, conf.SASLPassword)
	if err != nil {
		return nil, err
	}

	lastSchema, err := getLastSchema(conf.SchemaRegistry, conf.SASLUsername, conf.SASLPassword)
	if err != nil {
		return nil, err
	}

	return &auditTrails{
		serviceName: serviceName,
		producer:    producer,
		kafka:       conf,
		schema:      lastSchema,
		mysql:       mysql,
	}, nil
}

func initProducer(brokers string, username string, password string) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     username,
		"sasl.password":     password,
	})

	if err != nil {
		log.Printf("Failed to create audit trails producer: %s", err)
		return nil, err
	}

	log.Printf("Producer: %+v\n", producer)

	return producer, nil
}

// getLastSchema is used to get the last schema of audit trails
func getLastSchema(schemaRegistryURL string, user string, password string) (*srclient.Schema, error) {
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
	schemaRegistryClient.SetCredentials(user, password)

	schema, errSchema := schemaRegistryClient.GetLatestSchema(config.AUDIT_TRAILS_TOPIC_NAME + "-value")
	if errSchema != nil {
		log.Println("Failed to get schema: ", errSchema.Error())
		return nil, errSchema
	}

	log.Printf("getting schema: %+v\n", schema)

	return schema, nil
}

func (repo *auditTrails) Produce(data entity.AuditTrail) error {

	now := time.Now()
	data.Time = now.Format(constant.UNIX_TIME_LAYOUT)

	byteValue, err := json.Marshal(data)
	if err != nil {
		log.Println("Error: ", err)
		return err
	}

	schemaIDBytes := make([]byte, 4)
	schemaID := repo.schema.ID()
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schemaID))

	native, _, err := repo.schema.Codec().NativeFromTextual(byteValue)
	if err != nil {
		log.Println("error convert avro data from Textual")
		return err
	}

	valueBytes, err := repo.schema.Codec().BinaryFromNative(nil, native)
	if err != nil {
		log.Println("error append binary from Native")
		return err
	}

	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)

	key, _ := uuid.NewUUID()
	// log.Printf("message to send: %+v\n", string(recordValue))

	topicName := config.AUDIT_TRAILS_TOPIC_NAME

	kMessage := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topicName,
			Partition: kafka.PartitionAny,
		},
		Key: []byte(key.String()), Value: recordValue,
	}

	errProd := repo.producer.Produce(&kMessage, nil)
	if errProd != nil {
		log.Println("error producing message: ", errProd)
		return errProd
	}

	return nil
}

func (repo *auditTrails) Store(data entity.AuditTrail) error {
	err := repo.mysql.Create(data).Error
	if err != nil {
		log.Println("Query error storeAuditTrail when insert: ", err.Error())
		return err
	}

	return nil
}
