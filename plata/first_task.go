//go:build ignore

package main

/*
На входе есть отсортированные по возрастанию массивы чисел:
- suspects (подозреваемые)
- innocents (невиновные)

Необходимо из массива подозреваемых исключить всех невиновных.

Примеры:

Input: suspects = [1, 2, 3, 4, 5], innocents = [2, 4]
Output: [1, 3, 5]

Input: suspects = [3, 80, 123, 421, 936], innocents = [80, 936]
Output: [3, 123, 421]
*/

import "fmt"

func main() {
	suspects := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	innocents := []int{2, 5, 8}

	filtered := filter(suspects, innocents)
	fmt.Println(filtered)
}

func filter(suspects, innocents []int) []int {
	var res []int
	innocentsMap := make(map[int]struct{}, len(innocents))
	for _, i := range innocents {
		innocentsMap[i] = struct{}{}
	}

	for _, s := range suspects {
		if _, ok := innocentsMap[s]; ok {
			continue
		}
		res = append(res, s)
	}

	return res
}
