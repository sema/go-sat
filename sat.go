package main

import (
	"fmt"
)

type Formula []Clause

func (formula Formula) NumAtoms() uint64 {
	var atomIdMax uint64 = 0

	for _, clause := range formula {
		for _, literal := range clause {

			if uint64(literal.atom) > atomIdMax {
				atomIdMax = uint64(literal.atom)
			}
		}
	}

	return atomIdMax + 1
}

type Clause []Literal

type Literal struct {
	negated bool
	atom    Atom
}

type Atom uint64

type Solution []bool

type assignment struct {
	value    bool
	assigned bool
}

type stackAssignment struct {
	atom  Atom
	value bool
	branching bool
}

type assignments struct {
	raa          []assignment
	stackPointer uint64
	stackMax     uint64
	stack        []stackAssignment
}

func newAssignments(numAtoms uint64) *assignments {

	var assignments = new(assignments)

	assignments.raa = make([]assignment, numAtoms)
	assignments.stackPointer = 0
	assignments.stackMax = numAtoms
	assignments.stack = make([]stackAssignment, numAtoms)

	return assignments
}

func (assignments *assignments) Solution() Solution {
	var solution = make(Solution, len(assignments.raa))

	for i, assignment := range assignments.raa {
		solution[i] = assignment.value
	}

	return solution
}

func (assignments *assignments) HasUnassigned() bool {
	return assignments.stackPointer < assignments.stackMax
}

func (assignments *assignments) GetFirstUnassigned() (Atom, bool) {

	for i := range assignments.raa {
		if !assignments.raa[i].assigned {
			return Atom(i), true
		}
	}

	return Atom(0), false
}

func (assignments *assignments) PushAssignment(atom Atom, value bool, branching bool) bool {

	if uint64(atom) < assignments.stackMax {

		if assignments.raa[uint64(atom)].assigned {
			panic("Assignment to already assigned atom")
		}

		assignments.raa[uint64(atom)] = assignment{value, true}
		assignments.stack[assignments.stackPointer] = stackAssignment{atom, value, branching}
		assignments.stackPointer++
		return true
	}

	return false
}

func (assignments *assignments) PopAssignment() (stackAssignment, bool) {

	if assignments.stackPointer == 0 {
		return stackAssignment{}, false
	}

	assignments.stackPointer--

	var last = assignments.stack[assignments.stackPointer]

	assignments.raa[uint64(last.atom)].assigned = false
	return last, true
}

func (assignments *assignments) PartiallySatisfies(formula Formula) bool {

	for _, clause := range formula {

		var isSatisfied = false
		var hasUnassigned = false

		for _, literal := range clause {

			var assignment = assignments.raa[uint64(literal.atom)]

			hasUnassigned = hasUnassigned || !assignment.assigned

			isSatisfied =
				isSatisfied ||
					(literal.negated && !assignment.value) ||
					(!literal.negated && assignment.value)
		}

		if !isSatisfied && !hasUnassigned {
			return false
		}
	}

	return true

}

func (assignments *assignments) UnitProp(formula Formula) bool {

	dirty := true

	for dirty {

		dirty = false

		for _, clause := range formula {

			lastUnassigned := Literal{}
			numUnassigned := uint64(0)
			isSatisfied := false

			for _, literal := range clause {

				var assignment = assignments.raa[uint64(literal.atom)]

				if !assignment.assigned {
					numUnassigned++
					lastUnassigned = literal
				} else {

					isSatisfied =
						isSatisfied ||
						(literal.negated && !assignment.value) ||
						(!literal.negated && assignment.value)

				}

			}

			if numUnassigned == 1 && !isSatisfied {
				// if clause is a unit clause and not satisfied, select and set dirty

				assignments.PushAssignment(lastUnassigned.atom, !lastUnassigned.negated, false)
				dirty = true

			} else if (numUnassigned == 0 && !isSatisfied) {
				// if clause is not satisfied, set conflict
				
				return true
			}
		}
	}

	return false

}

func Solve(formula Formula) (Solution, bool) {

	var numAtoms = formula.NumAtoms()
	var assignments = newAssignments(numAtoms)

	var forceAlternative = false

	for {

		// success condition

		if !assignments.HasUnassigned() {
			return assignments.Solution(), true
		}

		// pick -- try new assignment

		var atom, ok = assignments.GetFirstUnassigned()
		if !ok {
			panic(fmt.Sprintf("Could not find an unassigned value, even though one should exist"))
		}

		assignments.PushAssignment(atom, forceAlternative, true)
		forceAlternative = false

		// deduce -- unit propagation

		var conflict = assignments.UnitProp(formula)

		// resolve - if we have a conflict

		conflict = conflict || !assignments.PartiallySatisfies(formula)

		if conflict {

			// backtrack

			for {

				var last, ok = assignments.PopAssignment()

				if !ok {
					// We can't backtrack further, so stop the search
					return nil, false
				}

				if last.value == false && last.branching {
					// We have backtracked to a false branch, stop and explore the true branch next iteration
					break
				}
			}

			forceAlternative = true
		}
	}

}
