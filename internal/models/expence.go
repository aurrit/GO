package models

type Category struct {
    ID   int
    Name string
}

type Expense struct {
    
    ID           int
    Amount       float64
    Description  string
    Date         string
    CategoryID   int
    CategoryName string
}
type ExpenseFilter struct {
    CategoryID int
    DateFrom   string
    DateTo     string
}