package mux

import (
	"net/http"
	"testing"
)

func TestEndpointUse(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf(err.(error).Error())
		}
	}()

	err := "Failed nodeImpl use."
	ep := &endpoint{}
	tVal := ""
	tVal2 := ""

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tVal = "A"
	})
	p := PluginFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tVal2 = "B"
		next(w, r)
	})

	ep.handler = handler
	ep.Use(p)

	r, _ := http.NewRequest("GET", "", nil)
	ep.chain.run(nil, r)
	if tVal != "A" {
		t.Errorf(err)
	}
	if tVal2 != "B" {
		t.Errorf(err)
	}
}

func TestEndpointServeHTTP(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf(err.(error).Error())
		}
	}()

	err := "Failed node ServeHTTP."
	ep := &endpoint{}
	tVal := ""
	tVal2 := ""

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tVal = "A"
	})
	p := PluginFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tVal2 = "C"
		next(w, r)
	})

	ep.handler = handler
	r, _ := http.NewRequest("GET", "", nil)
	ep.ServeHTTP(nil, r)
	if tVal != "A" {
		t.Errorf(err)
	}

	// Test run chain
	ep.Use(p)
	tVal = ""
	r, _ = http.NewRequest("GET", "", nil)
	ep.ServeHTTP(nil, r)
	if tVal != "A" {
		t.Errorf(err)
	}
	if tVal2 != "C" {
		t.Errorf(err)
	}
}
