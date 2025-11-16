package services_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"play.ground/generic-data-collector/internal/services"
)

// TestMockConsumer_WorkerSimulation tests the worker's core logic conceptually
// by simulating the message processing loop and capturing stdout.
func TestMockConsumer_WorkerSimulation(t *testing.T) {
	// 1. Setup MockConsumer and Subscription
	mockConsumer := services.NewMockConsumer()
	msgChan, err := mockConsumer.Subscribe("metrics")
	assert.NoError(t, err)

	// 2. Capture stdout
	old := os.Stdout // keep backup of original stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 3. Simulate Worker's Message Processing Loop
	// This goroutine mimics the logic in cmd/consumer/main.go
	done := make(chan struct{})
	go func() {
		for msg := range msgChan {
			// This is the exact logic from cmd/consumer/main.go
			fmt.Printf("Received message: %s\n", string(msg.Data()))
		}
		close(done)
	}()

	// 4. Inject a dummy message
	// dummyMessage := []byte(`{"value": 42}`)
	dummyMessage := services.NewNonAckPubSubMessage([]byte(`{"value": 42}`))
	mockConsumer.SendMessage(dummyMessage)

	// Give the goroutine a moment to process the message
	time.Sleep(10 * time.Millisecond)

	// 5. Stop capturing stdout and restore original
	w.Close()
	os.Stdout = old

	// 6. Read captured output
	out, _ := io.ReadAll(r)
	output := string(out)

	// 7. Assert the output
	expectedOutput := fmt.Sprintf("Received message: %s\n", string(dummyMessage.Data()))
	assert.Contains(t, output, expectedOutput, "Stdout should contain the processed message")

	// 8. Cleanup
	mockConsumer.Close()
	<-done // Wait for the worker goroutine to finish
}
