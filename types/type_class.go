// The MIT License (MIT)
//
// Copyright (c) 2019 West Damron
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package types

import (
	"github.com/wdamron/poly/internal/util"
)

// Parameterized type-class
type TypeClass struct {
	Name      string
	Param     Type
	Methods   map[string]*Arrow
	Super     map[string]*TypeClass
	Sub       map[string]*TypeClass
	Instances []*Instance
}

// Instance of a parameterized type-class
type Instance struct {
	TypeClass *TypeClass
	Param     Type
	Methods   map[string]*Arrow
}

// InstanceConstraint constrains a type-variable to types which implement a type-class.
type InstanceConstraint struct {
	TypeClass *TypeClass
}

// Create a new named/parameterized type-class with a set of method declarations.
func NewTypeClass(name string, param Type, methods map[string]*Arrow) *TypeClass {
	return &TypeClass{Name: name, Param: param, Methods: methods}
}

// Add a super-class to the type-class.
func (sub *TypeClass) AddSuperClass(super *TypeClass) { super.AddSubClass(sub) }

// Add a sub-class to the type-class.
func (super *TypeClass) AddSubClass(sub *TypeClass) {
	for _, tc := range super.Sub {
		if tc.Name == sub.Name {
			return
		}
	}
	if sub.Super == nil {
		sub.Super = make(map[string]*TypeClass)
	}
	if super.Sub == nil {
		super.Sub = make(map[string]*TypeClass)
	}
	sub.Super[super.Name] = super
	super.Sub[sub.Name] = sub
}

// Add an instance to the type-class with param as the type-parameter.
func (tc *TypeClass) AddInstance(param Type, methods map[string]*Arrow) *Instance {
	inst := &Instance{TypeClass: tc, Param: param, Methods: methods}
	tc.Instances = append(tc.Instances, inst)
	return inst
}

// Visit all instances for the type-class and all sub-classes. Sub-classes will be visited first.
func (tc *TypeClass) FindInstance(found func(*Instance) bool) bool {
	seen := util.NewDedupeMap()
	ok, _ := tc.findInstance(seen, found)
	seen.Release()
	return ok
}

func (tc *TypeClass) findInstance(seen map[string]bool, found func(*Instance) bool) (ok, shouldContinue bool) {
	if seen[tc.Name] {
		return false, true
	}
	seen[tc.Name] = true
	for _, sub := range tc.Sub {
		if ok, shouldContinue = sub.findInstance(seen, found); !shouldContinue {
			return ok, false
		}
	}
	for _, inst := range tc.Instances {
		if found(inst) {
			return true, false
		}
	}
	return false, true
}