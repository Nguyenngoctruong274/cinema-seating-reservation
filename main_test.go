package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

type mockCinemaService struct {
	saveCalled bool
	ShouldFail bool
}

func (m *mockCinemaService) SaveDataCinema(filename string) error {
	m.saveCalled = true
	if m.ShouldFail {
		return fmt.Errorf("mock error")
	}
	return nil
}

func TestPanicRecoverySaveDataCinema(t *testing.T) {
	mockService := &mockCinemaService{}

	defer func() {
		if re := recover(); re != nil {
			log.Printf("Recovered from panic: %v", re)
			_ = mockService.SaveDataCinema(stateFile)
		}
		if !mockService.saveCalled {
			t.Errorf("SaveDataCinema was not called on panic")
		}
	}()

	// Simulate panic
	panic("test panic")
}

// /////////////////////////////////////////////////////////////////////////////////////
func handleShutdownWithCallback(
	cinemaService *mockCinemaService,
	stateFile string,
	quit <-chan os.Signal,
	onExit func(int),
) {
	<-quit
	if err := cinemaService.SaveDataCinema(stateFile); err != nil {
		log.Printf("Failed to save state: %v", err)
	}
	onExit(0)
}

func TestHandleShutdownWithCallback(t *testing.T) {
	mockService := &mockCinemaService{}
	quit := make(chan os.Signal, 1)

	called := false
	exitCode := 0
	onExit := func(code int) {
		called = true
		exitCode = code
	}

	go handleShutdownWithCallback(mockService, "mockfile.json", quit, onExit)

	// Gửi tín hiệu
	quit <- syscall.SIGTERM
	time.Sleep(100 * time.Millisecond) // đợi goroutine chạy xong

	if !called || exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if !mockService.saveCalled {
		t.Errorf("Expected SaveDataCinema to be called")
	}
}
