package site_parsing

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(input string) *Data {
	if len(input) != 3 {
		fmt.Println("Invalid input. Must have 3 parts (Site ID, Task Type (0 (Monitor) or 1 (Checkout)) & Action (0 (Start) or 1 (Stop)))")
		return nil
	}

	_split := strings.Split(input, "")
	SiteID, err := strconv.Atoi(_split[0])
	if err != nil {
		fmt.Println("Invalid input. Site ID must be a number")
		return nil
	}

	TaskTypeInt, err := strconv.Atoi(_split[1])
	if err != nil {
		fmt.Println("Invalid input. Task Type must be a number (0 (Monitor), 1 (Checkout), or 2 (Both))")
		return nil
	}

	Action, err := strconv.Atoi(_split[2])
	if err != nil {
		fmt.Println("Invalid input. Action must be a number (0 (Start) or 1 (Stop))")
		return nil
	}

	SiteID = SiteID - 1

	return &Data{
		SiteID:   SiteID,
		TaskType: TaskTypeInt,
		Action:   Action,
	}

}
