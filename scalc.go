package main

import (
	"bufio"
	"fmt"
	"github.com/scylladb/go-set/iset"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	argsWithoutProc := os.Args[1:]
	if len(argsWithoutProc) < 4 {
		log.Fatal("bad input: expression must at have least 4 operators")
	}

	res := calc(argsWithoutProc)
	for _, r := range res {
		fmt.Println(r)
	}
}

// Expr struct holds order of operations, brackets and sets from the input
type Expr struct {
	set      *iset.Set
	operator string
}

// calc evaluates the expression and returns resulting set as []int
func calc(args []string) []int {
	var inputExpr []Expr
	// parse user input arguments:
	// read file names into iset.Set and save them into Expr.set
	// save brackets and operators as Expr.operators
	for _, arg := range args {
		if strings.HasSuffix(arg, ".txt") {
			set := file2set(arg)
			inputExpr = append(inputExpr, Expr{set: set})
		} else {
			inputExpr = append(inputExpr, Expr{operator: arg})
		}
	}

	start := 0
	end := 0

	// expression simplification loop, it finds the most inner expression, evaluates it and
	// repeats this process until there is only 1 element left in the inputExpr slice
	for {
		for i, token := range inputExpr {
			if token.operator == "[" {
				start = i
			} else if token.operator == "]" {
				end = i

				operator := inputExpr[start+1]
				sets := inputExpr[start+2 : end]

				var isets []*iset.Set
				for _, s := range sets {
					isets = append(isets, s.set)
				}

				set := calcExpr(operator.operator, isets)

				// build new slice with expression substituted with set from evaluated expression
				updatedTokens := inputExpr[:start]
				updatedTokens = append(updatedTokens, Expr{set: set})
				updatedTokens = append(updatedTokens, inputExpr[end+1:]...)
				inputExpr = updatedTokens
				break
			}
		}

		// when inputExpr is simplified to 1 Expr - this is our result
		if len(inputExpr) == 1 {
			// cast resulting set to []int and order it
			intRes := inputExpr[0].set.List()
			sort.Ints(intRes)

			return intRes
		}

	}
}

// calcExpr evaluates expression with only sets in it,
func calcExpr(operator string, sets []*iset.Set) *iset.Set {

	switch operator {

	case "DIF":
		if len(sets) < 2 {
			log.Fatalf("DIF command can only be performed on sets with 2 or more elements, got: %v", sets)
		}
		return iset.Difference(sets[0], sets[1:]...)

	case "INT":
		return iset.Intersection(sets...)

	case "SUM":
		return iset.Union(sets...)

	default:
		log.Fatalf("bad command input, got: %v", operator)
		return nil
	}
}

// file2set opens a file by its name, reads content into iset.Set, assumes that location of the files in cwd
func file2set(s string) *iset.Set {
	f, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	fileSet := iset.New()
	for scanner.Scan() {
		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		fileSet.Add(number)
	}
	return fileSet
}
