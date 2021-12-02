package tui

import (
	"fmt"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/rivo/tview"
)

var defaultHouse = houseModel.HouseDto{}

type CurrentHouse struct {
	*tview.TextView
	house houseModel.HouseDto
}

func NewCurrentHouse(house houseModel.HouseDto) *CurrentHouse {
	c := &CurrentHouse{
		TextView: tview.NewTextView(),
		house:    house,
	}

	c.setText()

	return c
}

func (c *CurrentHouse) enrich(flex *tview.Flex) {
	flex.AddItem(c, 3, 0, false)
}

func (c *CurrentHouse) setText() {
	if c.house == defaultHouse {
		c.SetText("Not house selected")
	} else {
		c.SetText(fmt.Sprintf("%s - %s\n%s", c.house.CountryCode, c.house.City, c.house.Name))
	}
}
