package codegen

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

const verbose = false

type Finder interface {
	FindMessageByName(protoPkg string, messageName string) *WasmMessage
	FindServiceByName(protoPkg string, messageName string) *WasmService
	AddMessageType(wasmName, protoPackage, goPackage string, message *WasmMessage)
	AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService)
	AddressingNameFromMessage(currentPkg string, message *WasmMessage) string
	GoPackageOption(service []*WasmService, message []*WasmMessage) (string, error)
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

func (s *SimpleFinder) GoPackageOption(service []*WasmService, message []*WasmMessage) (string, error) {
	pkg := ""
	for _, svc := range service {
		for sr, m := range s.service {
			if m == svc {
				part := strings.Split(sr.goPackage, ";")
				if len(part) != 2 {
					return "", fmt.Errorf("service %s: cannot understand go package option '%s'",
						svc.GetName(), sr.goPackage)
				}
				pkg = part[1]
			}
		}
	}
	if len(service) == 0 {
		for _, msg := range message {
			//log.Printf("%s?%s", msg.GetName(), msg.GetGoPackage())
			raw := msg.GetGoPackage()
			part := strings.Split(raw, ";")
			if len(part) == 2 {
				pkg = part[1]
			}
		}
	}
	if pkg == "" {
		panic("no package")
	}
	return pkg, nil
}

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

var versionRegexp = regexp.MustCompile(`^(v[0-9]+)(\..*$)`)

func (s *SimpleFinder) AddressingNameFromMessage(currentPkg string, message *WasmMessage) string {
	for candidate, m := range s.message {
		if m.GetFullName() == message.GetFullName() {
			if m.GetProtoPackage() == currentPkg {
				//log.Printf("\t\tcomparing pcakages %s vs %s (yay %s)", m.GetProtoPackage(), currentPkg, m.GetName())
				if verbose {
					log.Printf("! [simplefinder addressing name] found %s in package [%d] (current pkg was %s)", m.GetName(), len(m.GetField()), currentPkg)
				}
				return m.GetName()
			}
			if verbose {
				log.Printf("! [simplefinder addressing name] found %s, but in diff package [%d] (current pkg was %s)", m.GetFullName(), len(m.GetField()), currentPkg)
			}
			part := strings.Split(m.GetFullName(), ".")
			if len(part) > 2 {
				nonMsgPart := strings.Join(part[1:len(part)-1], ".")
				if nonMsgPart == currentPkg {
					nonParts := strings.Split(nonMsgPart, ".")
					if isVersion(nonParts[len(nonParts)-1]) {
						nonMsgPart = strings.Join(nonParts[:len(nonParts)-1], ".")
					}
					p_2 := nonMsgPart
					p_1 := part[len(part)-1]
					// YYY ies crucial
					return fmt.Sprintf("%s.%s", p_2, p_1)
				}
				result := strings.Join(part[len(part)-2:], ".")
				if versionRegexp.MatchString(result) {
					// move back one
					// YYY ies crucial
					pick := append([]string{part[len(part)-3]}, part[len(part)-1:]...)
					result = strings.Join(pick, ".")
				}
				//return strings.Join(part[len(part)-2:], ".")
				return result
			}
			if verbose {
				log.Printf("! [simplefinder addressing name] found %s, but in diff package and I cant understand the splitting of the name, giving up[%d] (current pkg was %s)", m.GetFullName(), len(m.GetField()), currentPkg)
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
	if protoPackage != "google.protobuf" {
		if verbose {
			log.Printf("new search for (%s,%s)---------", protoPackage, name)
		}
	}

	shortName := nameToJustServiceOrMessage(name)
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

func nameToJustServiceOrMessage(name string) string {
	name = strings.TrimSuffix(name, ".proto")
	part := strings.Split(name, ".")
	if len(part) == 1 {
		return name
	}
	return part[len(part)-1]
}

func (s *SimpleFinder) FindServiceByName(protoPackage string, name string) *WasmService {
	// sanity check
	if !strings.HasPrefix(name, "."+protoPackage) {
		panic(fmt.Sprintf("can't understand service/type structure: [.%s,%s]",
			protoPackage, name))
	}
	shortName := nameToJustServiceOrMessage(name)
	for candidate, svc := range s.service {
		if candidate.protoPackage == protoPackage {
			if candidate.wasmName == shortName {
				if verbose {
					log.Printf("! [simplefinder service] found %s", svc.GetWasmServiceName())
				}
				return svc
			} else {
				if verbose {
					log.Printf("- [simplefinder service]  package (%s) but not name %s vs %s xxxx %#v",
						protoPackage, candidate.wasmName, shortName, candidate)
				}
			}
		} else {
			if verbose {
				log.Printf("  [simplefinder service] missed %s versus [%s,%s]",
					candidate.String(), protoPackage, name)
			}
		}
	}
	return nil

}

func AddFileContentToFinder(f Finder, pr *descriptorpb.FileDescriptorProto, lang LanguageText) {
	for _, m := range pr.GetMessageType() {
		msg := NewWasmMessage(pr, m, lang, f)
		f.AddMessageType(m.GetName(), pr.GetPackage(), pr.GetOptions().GetGoPackage(), msg)
	}
	for _, s := range pr.GetService() {
		log.Printf("xxxx %s: %v", s.GetName(), isServiceMarkedParigot(s.Options.String()))
		svc := NewWasmService(pr, s, lang, f)
		log.Printf("xxxx-->> adding %s,%s,%s => %s", s.GetName(), pr.GetPackage(), pr.GetOptions().GetGoPackage(), svc.GetWasmServiceName())
		f.AddServiceType(s.GetName(), pr.GetPackage(), pr.GetOptions().GetGoPackage(), svc)
	}
}
