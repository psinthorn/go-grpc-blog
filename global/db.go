package global

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB mongo.Database

func connectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal("MongoDB connect error", err.Error())
	}

	DB = *client.Database(dbname)
}

// ตรวจวัดประสิทธิภาพของ application หากมีการเชื่อมต่อเข้ามากหรือน้อย ก็ให้เพิ่มหรือลด การเชื่อมต่อตามการใช้งาน
// NewDBContext จะคืนค่าเป็น Context.Background() กลับเมื่อมีการเชื่อมต่อไปยัง database เป็นไปตามเงื่อนไข
func NewDBContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d*performance/100)
}

func ConnectToTestDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal("MongoDB connect error", err.Error())
	}

	DB = *client.Database(dbname + "_test")

}
