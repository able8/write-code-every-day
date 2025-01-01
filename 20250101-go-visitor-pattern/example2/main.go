package main

import "fmt"

// https://blog.matthiasbruns.com/golang-visitor-pattern

type Visitor interface {
	VisitManager(manager *Manager)
	VisitEngineer(engineer *Engineer)
}

type Employee interface {
	Accept(visitor Visitor)
}

type Manager struct {
	Name   string
	Salary float64
}

func (m *Manager) Accept(visitor Visitor) {
	visitor.VisitManager(m)
}

type Engineer struct {
	Name   string
	Salary float64
}

func (e *Engineer) Accept(visitor Visitor) {
	visitor.VisitEngineer(e)
}

type SalaryCalculator struct {
	TotalSalary float64
}

func (s *SalaryCalculator) VisitManager(manager *Manager) {
	s.TotalSalary += manager.Salary
}

func (s *SalaryCalculator) VisitEngineer(engineer *Engineer) {
	s.TotalSalary += engineer.Salary
}

func main() {
	employees := []Employee{
		&Manager{Name: "John", Salary: 5000},
		&Engineer{Name: "Mary", Salary: 4000},
	}

	calculator := &SalaryCalculator{}
	for _, employee := range employees {
		employee.Accept(calculator)
	}

	fmt.Println("Total salary:", calculator.TotalSalary)
}
