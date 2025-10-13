package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// App struct
type App struct {
	ctx    context.Context
	timer  *time.Timer
	cancel chan struct{}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.cancel = make(chan struct{})
}

// shutdown is called at termination
func (a *App) shutdown(ctx context.Context) {
	a.CancelTimer() // Ensure any running timer is stopped
}

// StartTimer starts a timer for a power action
func (a *App) StartTimer(action string, timeValue int, timeUnit string) string {
	// Cancel any existing timer
	if a.timer != nil {
		a.timer.Stop()
	}

	var duration time.Duration
	switch timeUnit {
	case "seconds":
		duration = time.Duration(timeValue) * time.Second
	case "minutes":
		duration = time.Duration(timeValue) * time.Minute
	case "hours":
		duration = time.Duration(timeValue) * time.Hour
	default:
		return "Invalid time unit"
	}

	a.timer = time.NewTimer(duration)

	go func() {
		select {
		case <-a.timer.C:
			// Timer fired, execute the command
			var cmd *exec.Cmd
			switch action {
			case "shutdown":
				cmd = exec.Command("shutdown", "/s", "/t", "0")
			case "restart":
				cmd = exec.Command("shutdown", "/r", "/t", "0")
			case "sleep":
				// This command puts the computer to sleep
				cmd = exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
			default:
				fmt.Println("Invalid action")
				return
			}

			if err := cmd.Run(); err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		case <-a.cancel:
			// Timer was cancelled
			fmt.Println("Timer cancelled")
		}
	}()

	return fmt.Sprintf("Timer started. PC will %s in %d %s.", action, timeValue, timeUnit)
}

// CancelTimer stops the current timer
func (a *App) CancelTimer() {
	if a.timer != nil {
		if !a.timer.Stop() {
			// Timer already fired or was stopped, try to drain the channel
			select {
			case <-a.timer.C:
			default:
			}
		}
		// Signal cancellation to the goroutine
		close(a.cancel)
		a.cancel = make(chan struct{}) // Create a new channel for the next timer
		fmt.Println("Timer cancelled by user")
	}
}