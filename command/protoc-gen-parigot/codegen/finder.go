package codegen

import (
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"
	"log"
	"strings"
)

const verbose = false

type Finder interface {
	FindMessageByName(protoPkg string, messageName string) *WasmMessage
	FindServiceByName(protoPkg string, messageName string) *WasmService
	AddMessageType(wasmName, protoPackage, goPackage string, message *WasmMessage)
	AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService)
	AddressingNameFromMessage(currentPkg string, message *WasmMessage) string
	Service() []*WasmService
	Message() []*WasmMessage
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

var hack = make(map[string]int)

func (s *SimpleFinder) AddMessageType(name, protoPackage, goPackage string, message *WasmMessage) {
	rec := NewMessageRecord(name, protoPackage, goPackage)
	message.Collect()
	if verbose {
		if rec.protoPackage != "google.protobuf" && rec.goPackage != "google.golang.org/protobuf/types/descriptorpb)" {
			log.Printf("adding message type %s [%d]", rec.String(), len(message.GetField()))

		}
	}
	s.message[rec] = message
}

func (s *SimpleFinder) AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService) {
	rec := NewServiceRecord(wasmName, protoPackage, goPackage)
	service.Collect()
	if verbose {
		if rec.protoPackage != "google.protobuf" && rec.goPackage != "google.golang.org/protobuf/types/descriptorpb)" {
			log.Printf("adding service type %s", rec.String())
		}

	}
	s.service[rec] = service
}

func (s *SimpleFinder) Message() []*WasmMessage {
	result := []*WasmMessage{}
	for _, v := range s.message {
		result = append(result, v)
	}
	return result
}

func (s *SimpleFinder) Service() []*WasmService {
	result := []*WasmService{}
	for _, v := range s.service {
		result = append(result, v)
	}
	return result
}
func (s *SimpleFinder) AddressingNameFromMessage(currentPkg string, message *WasmMessage) string {
	for candidate, m := range s.message {
		if m.GetFullName() == message.GetFullName() {
			if m.GetProtoPackage() == currentPkg {
				if verbose {
					log.Printf("! [simplefinder addressing name] found %s [%d] (current pkg was %s)", m.GetName(), len(m.GetField()), currentPkg)
				}
				return m.GetName()
			}
			if verbose {
				log.Printf("! [simplefinder addressing name] found %s [%d] (current pkg was %s)", m.GetFullName(), len(m.GetField()), currentPkg)
			}
			return m.GetFullName()
		} else {
			if verbose {
				if candidate.protoPackage != "google.protobuf" && candidate.goPackage != "google.golang.org/protobuf/types/descriptorpb)" {
					log.Printf("  [simplefinder addressing name] missed %s versus %s",
						m.GetFullName(), message.GetFullName())
				}
			}

		}
	}
	return ""
}

func (s *SimpleFinder) FindMessageByName(protoPackage string, name string) *WasmMessage {
	// sanity check
	if !strings.HasPrefix(name, "."+protoPackage) && protoPackage != "" {
		//debug.PrintStack()
		//panic(fmt.Sprintf("can't understand message/type structure: [%s,%s]",
		//	protoPackage, name))
	}
	if protoPackage != "google.protobuf" {
		if verbose {
			log.Printf("new search for (%s,%s)---------", protoPackage, name)
		}
	}

	shortName := LastSegmentOfPackage(name)
	for candidate, m := range s.message {
		if candidate.protoPackage == protoPackage {
			if candidate.WasmName() == shortName {
				if verbose {
					log.Printf("! [simplefinder message] found %s", m.GetWasmMessageName())
				}
				return m
			} else {
				if verbose {
					if candidate.protoPackage != "google.protobuf" && candidate.goPackage != "google.golang.org/protobuf/types/descriptorpb)" {
						log.Printf("- [simplefinder message]  match package (%s) but not name %s vs %s",
							protoPackage, candidate.WasmName(), shortName)
					}

				}
			}
		} else {
			if verbose {
				if candidate.protoPackage != "google.protobuf" && candidate.goPackage != "google.golang.org/protobuf/types/descriptorpb)" {
					log.Printf("  [simplefinder message] missed %s versus [%s,%s]",
						candidate.String(), protoPackage, name)
				}
			}
		}
	}
	// we did not find it inside the package that contains the protoPackage, which is
	// the person who will be addressing it?  we are now going to try to make sure
	// we can get the message in the direct way
	for candidate, m := range s.message {
		if "."+candidate.protoPackage+"."+candidate.WasmName() == name {
			if verbose {
				log.Printf("! [simplefinder message] found %s [%d] in a different package than %s", name, len(m.GetField()), protoPackage)
			}
			return m
		}
	}
	if verbose {
		log.Printf("- [simplefinder message] missed completely on %s", name)
	}
	return nil
}

func (s *SimpleFinder) FindServiceByName(protoPackage string, name string) *WasmService {
	// sanity check
	if !strings.HasPrefix(name, "."+protoPackage) {
		panic(fmt.Sprintf("can't understand service/type structure: [%s,%s]",
			protoPackage, name))
	}
	shortName := LastSegmentOfPackage(name)
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
	return nil

}

func AddFileContentToFinder(f Finder, pr *descriptorpb.FileDescriptorProto,
	lang LanguageText) {
	for _, m := range pr.GetMessageType() {
		msg := NewWasmMessage(pr, m, lang, f)
		f.AddMessageType(m.GetName(), pr.GetPackage(), pr.GetOptions().GetGoPackage(), msg)
	}
	for _, s := range pr.GetService() {
		svc := NewWasmService(pr, s, lang, f)
		f.AddServiceType(s.GetName(), pr.GetPackage(), pr.GetOptions().GetGoPackage(), svc)
	}
}
