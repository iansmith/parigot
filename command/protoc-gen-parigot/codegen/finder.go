package codegen

import (
	"log"
	"strings"
)

const verbose = true

type Finder interface {
	FindMessageByName(protoPkg string, messageName string, next Finder) *WasmMessage
	FindServiceByName(protoPkg string, messageName string, next Finder) *WasmService
	AddMessageType(wasmName, protoPackage, goPackage string, message *WasmMessage)
	AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService)
}

type SimpleFinder struct {
	message map[*MessageRecord]*WasmMessage
	service map[*ServiceRecord]*WasmService
}

func NewSimpleFinder() *SimpleFinder {
	return &SimpleFinder{
		service: make(map[*ServiceRecord]*WasmService),
		message: make(map[*MessageRecord]*WasmMessage),
	}
}
func (s *SimpleFinder) AddMessageType(wasmName, protoPackage, goPackage string, message *WasmMessage) {
	rec := NewMessageRecord(wasmName, protoPackage, goPackage)
	s.message[rec] = message
}

func (s *SimpleFinder) AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService) {
	rec := NewServiceRecord(wasmName, protoPackage, goPackage)
	s.service[rec] = service
}

func (s *SimpleFinder) AllMessages() []*WasmMessage {
	result := []*WasmMessage{}
	for _, v := range s.message {
		result = append(result, v)
	}
	return result
}

func (s *SimpleFinder) FindMessageByName(protoPackage string, name string, next Finder) *WasmMessage {
	// sanity check
	if !strings.HasPrefix(name, "."+protoPackage) {
		log.Fatalf("can't understand message/type structure: [%s,%s]",
			protoPackage, name)
	}
	shortName := lastSegmentOfPackage(name)
	for candidate, m := range s.message {
		if candidate.protoPackage == protoPackage {
			if candidate.wasmName == shortName {
				if verbose {
					log.Printf("! [simplefinder message] found %s", m.GetWasmMessageName())
				}
				return m
			} else {
				if verbose {
					log.Printf("- [simplefinder message]  match package (%s) but not name %s vs %s",
						protoPackage, candidate.wasmName, shortName)
				}
			}
		} else {
			if verbose {
				log.Printf("  [simplefinder message] missed %s versus [%s,%s]",
					candidate.String(), protoPackage, name)
			}
		}
	}
	if next != nil {
		log.Printf("trying next finder...")
		next.FindMessageByName(protoPackage, name, nil)
	}
	return nil

}

func (s *SimpleFinder) FindServiceByName(protoPackage string, name string, next Finder) *WasmService {
	// sanity check
	if !strings.HasPrefix(name, "."+protoPackage) {
		log.Fatalf("can't understand service/type structure: [%s,%s]",
			protoPackage, name)
	}
	shortName := lastSegmentOfPackage(name)
	for candidate, svc := range s.service {
		if candidate.protoPackage == protoPackage {
			if candidate.wasmName == shortName {
				log.Printf("! [simplefinder service] found %s", svc.GetWasmServiceName())
				return svc
			} else {
				log.Printf("- [simplefinder service]  package (%s) but not name %s vs %s",
					protoPackage, candidate.wasmName, shortName)
			}
		} else {
			log.Printf("  [simplefinder service] missed %s versus [%s,%s]",
				candidate.String(), protoPackage, name)
		}
	}
	if next != nil {
		log.Printf("trying next finder...")
		next.FindMessageByName(protoPackage, name, nil)
	}
	return nil

}
