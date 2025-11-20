// Package rabbitmq provides examples of how to use the RabbitMQ package
package rabbitmq

// Example usage:
//
// func main() {
//     // Create a new RabbitMQ connection
//     conn, err := New("amqp://guest:guest@localhost:5672/",
//         WithReconnectInterval(10*time.Second),
//         WithHeartbeatInterval(30*time.Second),
//     )
//     if err != nil {
//         log.Fatal("Failed to connect to RabbitMQ:", err)
//     }
//     defer conn.Close()
//
//     // Create a channel
//     channel, err := conn.CreateChannel("example")
//     if err != nil {
//         log.Fatal("Failed to create channel:", err)
//     }
//
//     // Use the channel for publishing/consuming messages
//     // ...
//
//     // Get a channel by name
//     existingChannel, exists := conn.GetChannel("example")
//     if exists {
//         // Use existingChannel
//     }
// }
