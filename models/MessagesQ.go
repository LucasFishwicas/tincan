package models

// Defining a queue type and attaching enqueue() and dequeue() methods
type MessageQ struct {
    Messages []map[string]string
    Head int
    Tail int
    Length int
    Capacity int
}
func CreateQ(capacity int) *MessageQ {
    return &MessageQ{
        Messages: make([]map[string]string, capacity),
        Head: 0,
        Tail: 0,
        Length: 0,
        Capacity: capacity,
    }
}
func (self *MessageQ) Enqueue(user string, ipAddr string, message string) {
    if self.Length == self.Capacity {
        self.Dequeue()
    }
    messageMap := map[string]string{
        "user": user,
        "ipAddr": ipAddr,
        "message": message,
    }
    if self.Length != 0 {
        self.Tail = (self.Tail+1) % self.Capacity
    }
    self.Messages[self.Tail] = messageMap
    self.Length++
}
func (self *MessageQ) Dequeue() map[string]string {
    if len(self.Messages) == 0 {
        return nil
    }
    message := self.Messages[self.Head]
    self.Head = (self.Head+1) % self.Capacity
    self.Length--
    return message
}
// ---- |

