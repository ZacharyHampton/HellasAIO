package title

func AddCheckout() {
	checkouts += 1
	updateTitle()
}

func AddCart() {
	carts += 1
	updateTitle()
}

func AddFailure() {
	failures += 1
	updateTitle()
}
