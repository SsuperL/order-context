package resources

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func setUp() {
	fmt.Println("test start--------")
}

func tearDown() {
	os.Remove("gorm.db")
	fmt.Println("test end-------")
}
