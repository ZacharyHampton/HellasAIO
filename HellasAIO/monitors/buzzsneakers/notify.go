package buzzsneakersmonitor

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"time"
)

func notifyTasks(m *task.Task, b *BuzzSneakersInternal) task.TaskState {
	go func() {
		_ = m.NotifyTasks(&BuzzMonitorData{CheckoutData: m.CheckoutData, ProductId: b.ProductId, Products: b.Products})
	}()
	/*err := m.NotifyTasks(&BuzzMonitorData{CheckoutData: m.CheckoutData, ProductId: b.ProductId, Products: b.Products})
	if err != nil {
		logs.Log(m, "Error notifying checkout tasks.")
		time.Sleep(m.Delay)
		return NOTIFY_TASKS
	}*/

	logs.Log(m, fmt.Sprintf("Attempted to send notifications. Resending in %s seconds...", m.Delay))
	time.Sleep(m.Delay)
	return GET_PRODUCT_INFO
}
