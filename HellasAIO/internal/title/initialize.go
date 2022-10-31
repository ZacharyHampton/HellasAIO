package title

var (
	carts     int
	checkouts int
	failures  int
)

func Initialize() {
	carts = 0
	checkouts = 0
	failures = 0

	updateTitle()
}
