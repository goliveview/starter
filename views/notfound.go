package views

import (
	glv "github.com/goliveview/controller"
)

type NotfoundView struct {
	glv.DefaultView
}

func (n *NotfoundView) Content() string {
	return "./templates/404.html"
}

func (n *NotfoundView) Layout() string {
	return "./templates/layouts/error.html"
}
