package buzzsneakers

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	buzzsneakersmonitor "github.com/HellasAIO/HellasAIO/monitors/buzzsneakers"
)

func waitForMonitor(c *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	c.MonitorData = make(chan interface{})
	logs.Log(c, "Waiting for monitor...")

	for {
		select {
		case data := <-c.MonitorData:
			logs.Log(c, "Received notification. Checking out...")
			monitorData := data.(*buzzsneakersmonitor.BuzzMonitorData)
			c.CheckoutData = monitorData.CheckoutData
			b.Products = monitorData.Products
			b.ProductID = monitorData.ProductId

			if c.Size != "random" {
				c.CheckoutData.Size = c.Size
			}

			return ADD_TO_CART
		}
	}
}
