package main

import (
	"watn3y/steamsalty/botIO"
	"watn3y/steamsalty/steam"
)

func bot() {
	steam.StartWatchers(botIO.Authenticate())
}
