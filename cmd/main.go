package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

func checkDir(s string) bool {
	if strings.Contains(s, ".") || strings.Contains(s, "..") || strings.Contains(s, "/") || strings.Contains(s, "\\") || strings.Contains(s, "~") || strings.Contains(s, "*") {
		return false
	}

	return true
}

func printUsage() {
	fmt.Println("Coffee Shop Management System\n")
	fmt.Println("Usage:")
	fmt.Println(" hot-coffee [--port <N>] [--dir <S>]")
	fmt.Println(" hot-coffee --help\n")
	fmt.Println("Options:")
	fmt.Println(" --help     Show this screen.")
	fmt.Println(" --port N   Port number")
	fmt.Println(" --dir S    Path to the data directory")
}

func createDir(dir string) {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: dir does not created\n")
		os.Exit(1)
	}

	fd, err := os.OpenFile(dir+"/"+"order.json", os.O_RDWR|os.O_CREATE, 0o644)
	defer fd.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: order.json does not created\n")
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	for i := range args {
		if len(args[i]) == 6 && args[i] == "--help" || len(args[i]) == 2 && args[i] == "-h" {
			printUsage()
			os.Exit(0)
		}
	}

	dir := flag.String("dir", "data", "")
	port := flag.String("port", "3000", "")
	flag.Parse()

	number, _ := strconv.Atoi(*port)
	if number < 1024 || number > 49151 {
		fmt.Fprintf(os.Stderr, "Error: invalid port number\n")
		os.Exit(1)
	}
	if !checkDir(*dir) {
		fmt.Fprintf(os.Stderr, "Error: invalid dir name\n")
		os.Exit(1)
	}

	createDir(*dir)

	var dalInter dal.MenuDalInter = dal.ReturnMenuRepo(*dir + "/menu_items.json")
	var dalInv dal.InventDal = dal.ReturnInvent(*dir + "/inventory.json")
	var serInter service.MenuServiceInter = service.ReturnMenuService(dalInter, dalInv)
	var serInvent service.InventServiceInter = service.ReturnInvent(dalInv)

	menuHand := handler.ReturnMenuHand(serInter)
	inventhand := handler.ReturnInventHand(serInvent)

	o := dal.InitOrder("./" + *dir + "/order.json")
	orderService := &service.OrderService{OrderRepository: o, MenuRepository: dalInter, InventRepository: dalInv}
	orderHandler := handler.OrderHandler{OrderService: orderService}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", orderHandler.PostHandler)
	mux.HandleFunc("GET /orders", orderHandler.GetHandler)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetIdHandler)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.PutHandler)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteHandler)
	mux.HandleFunc("POST /orders/{id}/close", orderHandler.PostCloseHandler)
	mux.HandleFunc("GET /reports/total-sales", orderHandler.GetTotalSalesHandler)
	mux.HandleFunc("GET /reports/popular-items", orderHandler.GetPopularItemsHandler)

	mux.HandleFunc("GET /menu", menuHand.GetAllMenu)
	mux.HandleFunc("GET /menu/{ID}", menuHand.GetMenuId)
	mux.HandleFunc("DELETE /menu/{ID}", menuHand.DeleteMenuId)
	mux.HandleFunc("POST /menu", menuHand.PostMenu)
	mux.HandleFunc("PUT /menu/{ID}", menuHand.PutMenu)

	mux.HandleFunc("POST /inventory", inventhand.PostInvent)
	mux.HandleFunc("GET /inventory", inventhand.GetInvent)
	mux.HandleFunc("GET /inventory/{ID}", inventhand.GetInventId)
	mux.HandleFunc("PUT /inventory/{ID}", inventhand.PutInventId)
	mux.HandleFunc("DELETE /inventory/{ID}", inventhand.DeleteInventId)

	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
